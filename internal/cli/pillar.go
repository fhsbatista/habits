package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"habits/internal/usecase"
)

func NewPillarCommand(createPillar *usecase.CreatePillar, listPillars *usecase.ListPillars) func(args []string) {
	return func(args []string) {
		if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
			fmt.Println("Uso:")
			fmt.Println("  habits pillar add <nome>       Cria um novo pilar")
			fmt.Println("  habits pillar list             Lista todos os pilares")
			return
		}
		switch args[0] {
		case "add":
			pillarAdd(createPillar, args[1:])
		case "list":
			pillarList(listPillars)
		default:
			fmt.Fprintf(os.Stderr, "subcomando desconhecido: %s\n", args[0])
			os.Exit(1)
		}
	}
}

func pillarList(uc *usecase.ListPillars) {
	pillars, err := uc.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}
	if len(pillars) == 0 {
		fmt.Println("Nenhum pilar cadastrado.")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNome")
	fmt.Fprintln(w, "--\t----")
	for _, p := range pillars {
		fmt.Fprintf(w, "%d\t%s\n", p.ID, p.Name)
	}
	w.Flush()
}

func pillarAdd(uc *usecase.CreatePillar, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "uso: habits pillar add <nome>")
		os.Exit(1)
	}
	p, err := uc.Execute(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Pilar criado: [%d] %s\n", p.ID, p.Name)
}
