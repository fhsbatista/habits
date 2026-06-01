package cli

import (
	"fmt"
	"os"

	"habits/internal/domain"
	"habits/internal/usecase"
)

func NewStartCommand(checkIn *usecase.CheckIn) func(args []string) {
	return func(args []string) {
		if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
			fmt.Println("Uso:")
			fmt.Println("  habits start habit <id>   Inicia uma sessão para um hábito")
			fmt.Println("  habits start task <id>    Inicia uma sessão para uma tarefa")
			return
		}
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "uso: habits start <habit|task> <id>")
			os.Exit(1)
		}

		var refType domain.SessionRefType
		switch args[0] {
		case "habit":
			refType = domain.SessionRefHabit
		case "task":
			refType = domain.SessionRefTask
		default:
			fmt.Fprintf(os.Stderr, "tipo inválido %q: use habit ou task\n", args[0])
			os.Exit(1)
		}

		var id int64
		if _, err := fmt.Sscan(args[1], &id); err != nil {
			fmt.Fprintf(os.Stderr, "ID inválido: %s\n", args[1])
			os.Exit(1)
		}

		s, err := checkIn.Execute(refType, id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Sessão iniciada: %s %d às %s\n", args[0], id, s.StartedAt.Format("15:04:05"))
	}
}

func NewStopCommand(checkOut *usecase.CheckOut) func(args []string) {
	return func(args []string) {
		if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
			fmt.Println("Uso:")
			fmt.Println("  habits stop   Finaliza a sessão em andamento")
			return
		}

		s, err := checkOut.Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "erro: %v\n", err)
			os.Exit(1)
		}
		duration := s.FinishedAt.Sub(s.StartedAt)
		fmt.Printf("Sessão finalizada. Duração: %s\n", formatDuration(duration))
	}
}
