package cli

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

type command struct {
	fn   func(args []string)
	desc string
}

type App struct {
	cmds map[string]command
}

func NewApp() *App {
	return &App{cmds: map[string]command{}}
}

func (a *App) Register(name, desc string, fn func(args []string)) {
	a.cmds[name] = command{fn: fn, desc: desc}
}

func (a *App) Run(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" || args[0] == "help" {
		a.printHelp()
		return
	}

	cmd, ok := a.cmds[args[0]]
	if !ok {
		fmt.Fprintf(os.Stderr, "comando desconhecido: %s\n\n", args[0])
		a.printHelp()
		os.Exit(1)
	}
	cmd.fn(args[1:])
}

func (a *App) printHelp() {
	fmt.Println("habits — rastreador de hábitos e tarefas")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  habits <comando> [opções]")
	fmt.Println()
	fmt.Println("Comandos disponíveis:")

	names := make([]string, 0, len(a.cmds))
	for n := range a.cmds {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, n := range names {
		fmt.Printf("  %-18s %s\n", n, a.cmds[n].desc)
	}
	fmt.Println("\nUse 'habits <comando> --help' para detalhes de cada comando.")
}

func newFlagSet(name, usage string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: habits %s\n\n", usage)
		fs.PrintDefaults()
	}
	return fs
}
