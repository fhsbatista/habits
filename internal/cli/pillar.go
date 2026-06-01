package cli

import (
	"fmt"
	"os"

	"habits/internal/usecase"
)

func NewPillarCommand(createPillar *usecase.CreatePillar) func(args []string) {
	return func(args []string) {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "uso: habits pillar <add|list|remove>")
			os.Exit(1)
		}
		switch args[0] {
		case "add":
			pillarAdd(createPillar, args[1:])
		default:
			fmt.Fprintf(os.Stderr, "subcomando desconhecido: %s\n", args[0])
			os.Exit(1)
		}
	}
}

func pillarAdd(uc *usecase.CreatePillar, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "uso: habits pillar add <nome>")
		os.Exit(1)
	}
	name := args[0]

	p, err := uc.Execute(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Pilar criado: [%d] %s\n", p.ID, p.Name)
}
