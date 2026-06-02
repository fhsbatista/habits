package cli

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/term"

	"habits/internal/domain"
	"habits/internal/usecase"
)

const (
	titleWidth  = 22
	minBarWidth = 20
	ansiReset   = "\033[0m"
	ansiGreen   = "\033[32m"
	ansiBlue    = "\033[34m"
	ansiYellow  = "\033[33m"
	ansiRed     = "\033[31m"
	ansiOrange  = "\033[38;5;208m"
	ansiMagenta = "\033[35m"
)

func NewLogCommand(dailyTimeline *usecase.DailyTimeline, habitGrid *usecase.HabitGrid) func(args []string) {
	return func(args []string) {
		fs := newFlagSet("log", "log [--date YYYY-MM-DD] [--days N] [--month YYYY-MM]")
		dateStr := fs.String("date", "", "data no formato YYYY-MM-DD (padrão: hoje)")
		daysN := fs.Int("days", 0, "exibe os últimos N dias no grid")
		monthStr := fs.String("month", "", "mês para o grid no formato YYYY-MM")
		fs.Parse(args)

		day := time.Now()
		if *dateStr != "" {
			var err error
			day, err = time.ParseInLocation("2006-01-02", *dateStr, time.Local)
			if err != nil {
				fmt.Fprintf(os.Stderr, "data inválida %q: use o formato YYYY-MM-DD\n", *dateStr)
				os.Exit(1)
			}
		}

		entries, err := dailyTimeline.Execute(day)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro: %v\n", err)
			os.Exit(1)
		}

		termWidth := terminalWidth()
		barWidth := termWidth - titleWidth - 1
		if barWidth < minBarWidth {
			barWidth = minBarWidth
		}
		minsPerBlock := 1440.0 / float64(barWidth)

		printTimelineHeader(barWidth, minsPerBlock)

		if len(entries) == 0 {
			fmt.Println("  (nenhuma sessão registrada)")
		} else {
			for _, e := range entries {
				printTimelineEntry(e, barWidth, minsPerBlock, day)
			}
		}

		gridStart, gridEnd := gridRange(day, *daysN, *monthStr)
		grid, err := habitGrid.Execute(gridStart, gridEnd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro ao carregar grid: %v\n", err)
			os.Exit(1)
		}
		if len(grid.Habits) > 0 {
			fmt.Println()
			printHabitGrid(grid)
		}
	}
}

func gridRange(day time.Time, daysN int, monthStr string) (time.Time, time.Time) {
	if daysN > 0 {
		end := time.Date(day.Year(), day.Month(), day.Day()+1, 0, 0, 0, 0, day.Location())
		start := end.AddDate(0, 0, -daysN)
		return start, end
	}
	if monthStr != "" {
		t, err := time.ParseInLocation("2006-01", monthStr, time.Local)
		if err == nil {
			start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
			end := start.AddDate(0, 1, 0)
			return start, end
		}
	}
	end := time.Date(day.Year(), day.Month(), day.Day()+1, 0, 0, 0, 0, day.Location())
	start := end.AddDate(0, 0, -7)
	return start, end
}

func printTimelineHeader(barWidth int, minsPerBlock float64) {
	hours := []int{0, 6, 12, 18, 24}
	header := make([]byte, barWidth)
	for i := range header {
		header[i] = ' '
	}
	for _, h := range hours {
		pos := int(float64(h*60) / minsPerBlock)
		if pos >= barWidth {
			pos = barWidth - 1
		}
		label := fmt.Sprintf("%02d", h)
		for i, c := range label {
			if pos+i < barWidth {
				header[pos+i] = byte(c)
			}
		}
	}
	scale := int(minsPerBlock)
	if scale < 1 {
		scale = 1
	}
	fmt.Printf("%s  escala: 1█ = %dmin\n", string(header), scale)

	ruler := make([]byte, barWidth)
	for i := range ruler {
		ruler[i] = ' '
	}
	for _, h := range hours {
		pos := int(float64(h*60) / minsPerBlock)
		if pos >= barWidth {
			pos = barWidth - 1
		}
		ruler[pos] = '|'
	}
	fmt.Println(string(ruler))
}

func printTimelineEntry(e usecase.TimelineEntry, barWidth int, minsPerBlock float64, day time.Time) {
	localStart := e.Session.StartedAt.Local()
	localEnd := e.EffectiveEnd.Local()
	startMin := float64(localStart.Hour()*60 + localStart.Minute())
	endMin := float64(localEnd.Hour()*60 + localEnd.Minute())

	startPos := int(startMin / minsPerBlock)
	endPos := int(endMin / minsPerBlock)
	if endPos <= startPos {
		endPos = startPos + 1
	}
	if endPos > barWidth {
		endPos = barWidth
	}

	color := colorFor(e)
	block := color + strings.Repeat("█", endPos-startPos) + ansiReset

	line := make([]byte, barWidth)
	for i := range line {
		line[i] = ' '
	}

	prefix := string(line[:startPos])
	suffix := string(line[endPos:])

	name := e.Name
	if e.IsOpen {
		name += " (em andamento)"
	}
	if len([]rune(name)) > titleWidth-2 {
		name = string([]rune(name)[:titleWidth-5]) + "..."
	}

	fmt.Printf("%s%s%s  %s\n", prefix, block, suffix, name)
}

func colorFor(e usecase.TimelineEntry) string {
	if e.RefType == domain.SessionRefTask {
		return ansiMagenta
	}
	return habitColor(e.Color)
}

func habitColor(c domain.Color) string {
	switch c {
	case domain.ColorVerde:
		return ansiGreen
	case domain.ColorAzul:
		return ansiBlue
	case domain.ColorAmarelo:
		return ansiYellow
	case domain.ColorVermelho:
		return ansiRed
	case domain.ColorLaranja:
		return ansiOrange
	default:
		return ansiGreen
	}
}

func printHabitGrid(grid usecase.HabitGridResult) {
	const cellWidth = 3

	maxNameLen := 0
	for _, h := range grid.Habits {
		n := utf8.RuneCountInString(h.Name)
		if n > maxNameLen {
			maxNameLen = n
		}
	}
	if maxNameLen > titleWidth {
		maxNameLen = titleWidth
	}

	totalDays := len(grid.Days)
	percentages := make(map[int64]float64, len(grid.Habits))
	minPct, maxPct := 101.0, -1.0
	for _, h := range grid.Habits {
		count := 0
		for _, d := range grid.Days {
			if grid.Performed[h.ID][d] {
				count++
			}
		}
		var pct float64
		if totalDays > 0 {
			pct = float64(count) / float64(totalDays) * 100
		}
		percentages[h.ID] = pct
		if pct < minPct {
			minPct = pct
		}
		if pct > maxPct {
			maxPct = pct
		}
	}

	// header row: day numbers
	fmt.Printf("%*s", maxNameLen+1, "")
	for _, d := range grid.Days {
		fmt.Printf(" %02d ", d.Day())
	}
	fmt.Println()

	for _, h := range grid.Habits {
		name := h.Name
		if utf8.RuneCountInString(name) > maxNameLen {
			runes := []rune(name)
			name = string(runes[:maxNameLen-1]) + "…"
		}
		fmt.Printf("%-*s", maxNameLen+1, name)

		color := habitColor(h.Color)
		for _, d := range grid.Days {
			if grid.Performed[h.ID][d] {
				fmt.Printf(" %s■%s  ", color, ansiReset)
			} else {
				fmt.Printf("%*s", cellWidth+1, "")
			}
		}

		pct := percentages[h.ID]
		pctColor := percentageColor(pct, minPct, maxPct)
		fmt.Printf("  %s%3.0f%%%s\n", pctColor, pct, ansiReset)
	}
}

func percentageColor(pct, minPct, maxPct float64) string {
	var t float64
	if maxPct > minPct {
		t = (pct - minPct) / (maxPct - minPct)
	} else {
		t = 1.0
	}
	var r, g int
	if t <= 0.5 {
		r = 255
		g = int(t * 2 * 255)
	} else {
		r = int((1 - t) * 2 * 255)
		g = 255
	}
	return fmt.Sprintf("\033[38;2;%d;%d;0m", r, g)
}

func terminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}
