package cli

import (
	"flag"
	"fmt"
	"os"
)

type App struct {
	cmds map[string]func(args []string)
}

func NewApp() *App {
	return &App{cmds: map[string]func(args []string){}}
}

func (a *App) Register(name string, fn func(args []string)) {
	a.cmds[name] = fn
}

func (a *App) Run(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "uso: habits <comando> [opções]")
		os.Exit(1)
	}
	fn, ok := a.cmds[args[0]]
	if !ok {
		fmt.Fprintf(os.Stderr, "comando desconhecido: %s\n", args[0])
		os.Exit(1)
	}
	fn(args[1:])
}

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	return fs
}
