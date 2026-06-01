package usecase

import (
	"fmt"

	"habits/internal/domain"
)

type RemoveHabit struct {
	habitRepo domain.HabitRepository
}

func NewRemoveHabit(habitRepo domain.HabitRepository) *RemoveHabit {
	return &RemoveHabit{habitRepo: habitRepo}
}

func (uc *RemoveHabit) Execute(id int64) error {
	_, found, err := uc.habitRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("buscar hábito: %w", err)
	}
	if !found {
		return fmt.Errorf("hábito %d não encontrado", id)
	}
	if err := uc.habitRepo.Delete(id); err != nil {
		return fmt.Errorf("remover hábito: %w", err)
	}
	return nil
}
