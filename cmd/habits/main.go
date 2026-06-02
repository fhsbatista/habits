package main

import (
	"fmt"
	"os"

	"habits/internal/cli"
	"habits/internal/repository/sqlite"
	"habits/internal/usecase"
)

func main() {
	dbPath := os.Getenv("HABITS_DB")
	if dbPath == "" {
		home, _ := os.UserHomeDir()
		dbPath = home + "/.habits.db"
	}

	db, err := sqlite.Open(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro ao abrir banco: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	habitRepo := sqlite.NewHabitRepository(db)
	pillarRepo := sqlite.NewPillarRepository(db)
	taskRepo := sqlite.NewTaskRepository(db)
	sessionRepo := sqlite.NewSessionRepository(db)

	createHabit := usecase.NewCreateHabit(pillarRepo, habitRepo)
	createPillar := usecase.NewCreatePillar(pillarRepo)
	listPillars := usecase.NewListPillars(pillarRepo)
	listHabits := usecase.NewListHabits(habitRepo)
	removeHabit := usecase.NewRemoveHabit(habitRepo)
	dailyOverview := usecase.NewDailyOverview(pillarRepo, habitRepo, taskRepo)
	checkIn := usecase.NewCheckIn(habitRepo, taskRepo, sessionRepo)
	checkOut := usecase.NewCheckOut(sessionRepo)
	completeHabit := usecase.NewCompleteHabit(habitRepo, sessionRepo)
	dailyTimeline := usecase.NewDailyTimeline(habitRepo, taskRepo, sessionRepo)
	habitGrid := usecase.NewHabitGrid(habitRepo, sessionRepo)

	app := cli.NewApp()
	app.Register("pillar", "Gerencia pilares (add, list)", cli.NewPillarCommand(createPillar, listPillars))
	app.Register("habit", "Gerencia hábitos (add, list, remove)", cli.NewHabitCommand(createHabit, listHabits, removeHabit))
	app.Register("list", "Lista hábitos pendentes hoje e tarefas em andamento", cli.NewListCommand(dailyOverview))
	app.Register("ls", "Alias para list", cli.NewListCommand(dailyOverview))
	app.Register("complete", "Marca um hábito como realizado agora", cli.NewCompleteCommand(completeHabit))
	app.Register("start", "Inicia uma sessão para um hábito ou tarefa", cli.NewStartCommand(checkIn))
	app.Register("stop", "Finaliza a sessão em andamento", cli.NewStopCommand(checkOut))
	app.Register("log", "Exibe a timeline do dia e grid de hábitos", cli.NewLogCommand(dailyTimeline, habitGrid))

	app.Run(os.Args[1:])
}
