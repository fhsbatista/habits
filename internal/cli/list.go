package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"habits/internal/usecase"
)

func NewListCommand(dailyOverview *usecase.DailyOverview) func(args []string) {
	return func(args []string) {
		if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
			fmt.Println("Uso:")
			fmt.Println("  habits list   Lista hábitos pendentes hoje e tarefas em andamento")
			fmt.Println("  habits ls     Alias para list")
			return
		}
		result, err := dailyOverview.Execute(time.Now())
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro: %v\n", err)
			os.Exit(1)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		if len(result.DueHabits) > 0 {
			fmt.Fprintln(w, "[ HÁBITOS DE HOJE ]")
			fmt.Fprintln(w, "ID\tNome\tPilar")
			fmt.Fprintln(w, "--\t----\t-----")
			for _, ph := range result.DueHabits {
				if ph.Performed {
					fmt.Fprintf(w, "\033[9m%d\t%s\t%s\033[0m\n", ph.Habit.ID, ph.Habit.Name, ph.Pillar.Name)
				} else {
					fmt.Fprintf(w, "%d\t%s\t%s\n", ph.Habit.ID, ph.Habit.Name, ph.Pillar.Name)
				}
			}
			fmt.Fprintln(w)
		}

		fmt.Fprintln(w, "[ EXECUTAR ]")
		if len(result.ExecuteTasks) == 0 {
			fmt.Fprintln(w, "  (nenhuma tarefa)")
		} else {
			fmt.Fprintln(w, "ID\tTítulo\tDescrição\tSem update")
			fmt.Fprintln(w, "--\t------\t---------\t----------")
			for _, t := range result.ExecuteTasks {
				desc := t.Description
				if desc == "" {
					desc = "-"
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, t.Title, truncate(desc, 30), formatDuration(t.TimeSinceUpdate(time.Now())))
			}
		}
		fmt.Fprintln(w)

		fmt.Fprintln(w, "[ AGUARDAR ]")
		if len(result.WaitTasks) == 0 {
			fmt.Fprintln(w, "  (nenhuma tarefa)")
		} else {
			fmt.Fprintln(w, "ID\tTítulo\tDescrição\tSem update")
			fmt.Fprintln(w, "--\t------\t---------\t----------")
			for _, t := range result.WaitTasks {
				desc := t.Description
				if desc == "" {
					desc = "-"
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, t.Title, truncate(desc, 30), formatDuration(t.TimeSinceUpdate(time.Now())))
			}
		}

		w.Flush()
	}
}

func formatDuration(d time.Duration) string {
	switch {
	case d < time.Hour:
		m := int(d.Minutes())
		if m <= 1 {
			return "1 minuto"
		}
		return fmt.Sprintf("%d minutos", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hora"
		}
		return fmt.Sprintf("%d horas", h)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 dia"
		}
		return fmt.Sprintf("%d dias", days)
	}
}

func truncate(s string, max int) string {
	if len([]rune(s)) <= max {
		return s
	}
	return string([]rune(s)[:max-3]) + "..."
}

