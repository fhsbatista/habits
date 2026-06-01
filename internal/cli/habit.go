package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"habits/internal/domain"
	"habits/internal/usecase"
)

func NewHabitCommand(createHabit *usecase.CreateHabit, listHabits *usecase.ListHabits) func(args []string) {
	return func(args []string) {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "uso: habits habit <add|list>")
			os.Exit(1)
		}
		switch args[0] {
		case "add":
			habitAdd(createHabit, args[1:])
		case "list":
			habitList(listHabits, args[1:])
		default:
			fmt.Fprintf(os.Stderr, "subcomando desconhecido: %s\n", args[0])
			os.Exit(1)
		}
	}
}

func habitList(uc *usecase.ListHabits, args []string) {
	fs := newFlagSet("habit list")
	pillarID := fs.Int64("pillar", 0, "filtrar por pilar")
	fs.Parse(args)

	habits, err := uc.Execute(*pillarID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}
	if len(habits) == 0 {
		fmt.Println("Nenhum hábito cadastrado.")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNome\tCor\tDias")
	fmt.Fprintln(w, "--\t----\t---\t----")
	for _, h := range habits {
		days := make([]string, len(h.Frequency))
		for i, d := range h.Frequency {
			days[i] = string(d)
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", h.ID, h.Name, string(h.Color), strings.Join(days, ","))
	}
	w.Flush()
}

func habitAdd(uc *usecase.CreateHabit, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "uso: habits habit add <nome> --pillar <id> --days <dias>")
		os.Exit(1)
	}
	name := args[0]

	fs := newFlagSet("habit add")
	pillarID := fs.Int64("pillar", 0, "ID do pilar")
	days := fs.String("days", "", "dias da semana (ex: seg,qua,sex ou todos)")
	fs.Parse(args[1:])

	if *pillarID == 0 {
		fmt.Fprintln(os.Stderr, "erro: --pillar é obrigatório")
		os.Exit(1)
	}
	if *days == "" {
		fmt.Fprintln(os.Stderr, "erro: --days é obrigatório")
		os.Exit(1)
	}
	frequency := parseWeekdays(*days)
	color := promptColor()

	h, err := uc.Execute(usecase.CreateHabitInput{
		Name:      name,
		PillarID:  *pillarID,
		Color:     color,
		Frequency: frequency,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Hábito criado: [%d] %s\n", h.ID, h.Name)
}

func parseWeekdays(s string) []domain.Weekday {
	if s == "todos" {
		return []domain.Weekday{
			domain.WeekdayDom, domain.WeekdaySeg, domain.WeekdayTer,
			domain.WeekdayQua, domain.WeekdayQui, domain.WeekdaySex, domain.WeekdaySab,
		}
	}
	parts := strings.Split(s, ",")
	days := make([]domain.Weekday, 0, len(parts))
	for _, p := range parts {
		days = append(days, domain.Weekday(strings.TrimSpace(p)))
	}
	return days
}

func promptColor() domain.Color {
	validColors := []string{"azul", "verde", "laranja", "vermelho", "amarelo"}
	fmt.Printf("Escolha uma cor [%s] (Enter para aleatório): ", strings.Join(validColors, ", "))

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if input == "" {
		return domain.ColorAleatorio
	}
	for _, c := range validColors {
		if strings.EqualFold(input, c) {
			return domain.Color(c)
		}
	}
	fmt.Fprintf(os.Stderr, "cor inválida %q, usando aleatório\n", input)
	return domain.ColorAleatorio
}
