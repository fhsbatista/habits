package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"habits/internal/domain"
	"habits/internal/usecase"
)

const (
	titleWidth  = 22
	minBarWidth = 20
	ansiReset   = "\033[0m"
	ansiGreen   = "\033[32m"
	ansiMagenta = "\033[35m"
)

func NewLogCommand(dailyTimeline *usecase.DailyTimeline) func(args []string) {
	return func(args []string) {
		fs := newFlagSet("log")
		dateStr := fs.String("date", "", "data no formato YYYY-MM-DD")
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
			return
		}

		for _, e := range entries {
			printTimelineEntry(e, barWidth, minsPerBlock, day)
		}
	}
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
	startMin := float64(e.Session.StartedAt.Hour()*60 + e.Session.StartedAt.Minute())
	endMin := float64(e.EffectiveEnd.Hour()*60 + e.EffectiveEnd.Minute())

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
	return ansiGreen
}

func terminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}
