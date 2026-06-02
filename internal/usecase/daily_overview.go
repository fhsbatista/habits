package usecase

import (
	"fmt"
	"sort"
	"time"

	"habits/internal/domain"
)

type DueHabit struct {
	Habit     domain.Habit
	Pillar    domain.Pillar
	Performed bool
}

type DailyOverviewResult struct {
	DueHabits    []DueHabit
	ExecuteTasks []domain.Task
	WaitTasks    []domain.Task
}

type DailyOverview struct {
	pillarRepo  domain.PillarRepository
	habitRepo   domain.HabitRepository
	taskRepo    domain.TaskRepository
	sessionRepo domain.SessionRepository
}

func NewDailyOverview(pillarRepo domain.PillarRepository, habitRepo domain.HabitRepository, taskRepo domain.TaskRepository, sessionRepo domain.SessionRepository) *DailyOverview {
	return &DailyOverview{pillarRepo: pillarRepo, habitRepo: habitRepo, taskRepo: taskRepo, sessionRepo: sessionRepo}
}

func (uc *DailyOverview) Execute(now time.Time) (DailyOverviewResult, error) {
	pending, err := uc.pendingHabits(now)
	if err != nil {
		return DailyOverviewResult{}, fmt.Errorf("hábitos pendentes: %w", err)
	}

	execute, wait, err := uc.groupedTasks(now)
	if err != nil {
		return DailyOverviewResult{}, fmt.Errorf("tarefas: %w", err)
	}

	return DailyOverviewResult{
		DueHabits:    pending,
		ExecuteTasks: execute,
		WaitTasks:    wait,
	}, nil
}

func (uc *DailyOverview) pendingHabits(now time.Time) ([]DueHabit, error) {
	habits, err := uc.habitRepo.FindAll()
	if err != nil {
		return nil, err
	}

	sessions, err := uc.sessionRepo.FindByDay(now)
	if err != nil {
		return nil, err
	}

	var due []DueHabit
	for _, h := range habits {
		if !h.IsDueOn(domain.WeekdayFromTime(now)) {
			continue
		}
		pillar, _, err := uc.pillarRepo.FindByID(h.PillarID)
		if err != nil {
			return nil, err
		}
		due = append(due, DueHabit{Habit: h, Pillar: pillar, Performed: h.WasPerformedOn(now, sessions)})
	}
	return due, nil
}

func (uc *DailyOverview) groupedTasks(now time.Time) (execute, wait []domain.Task, err error) {
	tasks, err := uc.taskRepo.FindInProgress()
	if err != nil {
		return nil, nil, err
	}

	for _, t := range tasks {
		if t.NextAction != nil && t.NextAction.Type == domain.NextActionAguardar {
			wait = append(wait, t)
		} else {
			execute = append(execute, t)
		}
	}

	sortByTimeSinceUpdate(execute, now)
	sortByTimeSinceUpdate(wait, now)
	return execute, wait, nil
}

func sortByTimeSinceUpdate(tasks []domain.Task, now time.Time) {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].TimeSinceUpdate(now) > tasks[j].TimeSinceUpdate(now)
	})
}
