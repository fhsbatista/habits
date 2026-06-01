package usecase

import (
	"fmt"
	"sort"
	"time"

	"habits/internal/domain"
)

type PendingHabit struct {
	Habit  domain.Habit
	Pillar domain.Pillar
}

type DailyOverviewResult struct {
	PendingHabits []PendingHabit
	ExecuteTasks  []domain.Task
	WaitTasks     []domain.Task
}

type DailyOverview struct {
	pillarRepo domain.PillarRepository
	habitRepo  domain.HabitRepository
	taskRepo   domain.TaskRepository
}

func NewDailyOverview(pillarRepo domain.PillarRepository, habitRepo domain.HabitRepository, taskRepo domain.TaskRepository) *DailyOverview {
	return &DailyOverview{pillarRepo: pillarRepo, habitRepo: habitRepo, taskRepo: taskRepo}
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
		PendingHabits: pending,
		ExecuteTasks:  execute,
		WaitTasks:     wait,
	}, nil
}

func (uc *DailyOverview) pendingHabits(now time.Time) ([]PendingHabit, error) {
	habits, err := uc.habitRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var pending []PendingHabit
	for _, h := range habits {
		if !h.IsDueOn(domain.WeekdayFromTime(now)) {
			continue
		}
		pillar, _, err := uc.pillarRepo.FindByID(h.PillarID)
		if err != nil {
			return nil, err
		}
		pending = append(pending, PendingHabit{Habit: h, Pillar: pillar})
	}
	return pending, nil
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
