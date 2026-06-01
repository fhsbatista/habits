package usecase

import (
	"fmt"
	"time"

	"habits/internal/domain"
)

type HabitGridResult struct {
	Habits    []domain.Habit
	Days      []time.Time
	Performed map[int64]map[time.Time]bool
}

type HabitGrid struct {
	habitRepo   domain.HabitRepository
	sessionRepo domain.SessionRepository
}

func NewHabitGrid(habitRepo domain.HabitRepository, sessionRepo domain.SessionRepository) *HabitGrid {
	return &HabitGrid{habitRepo: habitRepo, sessionRepo: sessionRepo}
}

func (uc *HabitGrid) Execute(start, end time.Time) (HabitGridResult, error) {
	habits, err := uc.habitRepo.FindAll()
	if err != nil {
		return HabitGridResult{}, fmt.Errorf("buscar hábitos: %w", err)
	}

	sessions, err := uc.sessionRepo.FindByDateRange(start, end)
	if err != nil {
		return HabitGridResult{}, fmt.Errorf("buscar sessões: %w", err)
	}

	var days []time.Time
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		days = append(days, time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()))
	}

	performed := make(map[int64]map[time.Time]bool)
	for _, h := range habits {
		performed[h.ID] = make(map[time.Time]bool)
		for _, day := range days {
			if h.WasPerformedOn(day, sessions) {
				performed[h.ID][day] = true
			}
		}
	}

	return HabitGridResult{
		Habits:    habits,
		Days:      days,
		Performed: performed,
	}, nil
}
