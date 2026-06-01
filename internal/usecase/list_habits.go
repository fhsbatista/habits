package usecase

import (
	"fmt"

	"habits/internal/domain"
)

type ListHabits struct {
	habitRepo domain.HabitRepository
}

func NewListHabits(habitRepo domain.HabitRepository) *ListHabits {
	return &ListHabits{habitRepo: habitRepo}
}

func (uc *ListHabits) Execute(pillarID int64) ([]domain.Habit, error) {
	if pillarID != 0 {
		habits, err := uc.habitRepo.FindByPillar(pillarID)
		if err != nil {
			return nil, fmt.Errorf("listar hábitos por pilar: %w", err)
		}
		return habits, nil
	}
	habits, err := uc.habitRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("listar hábitos: %w", err)
	}
	return habits, nil
}
