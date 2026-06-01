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

	app := cli.NewApp()
	app.Register("pillar", cli.NewPillarCommand(createPillar, listPillars))
	app.Register("habit", cli.NewHabitCommand(createHabit, listHabits, removeHabit))
	app.Register("list", cli.NewListCommand(dailyOverview))
	app.Register("ls", cli.NewListCommand(dailyOverview))
	app.Register("start", cli.NewStartCommand(checkIn))

	app.Run(os.Args[1:])
}
