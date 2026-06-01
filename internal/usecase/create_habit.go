package usecase

import (
	"fmt"

	"habits/internal/domain"
)

type CreateHabitInput struct {
	Name      string
	PillarID  int64
	Color     domain.Color
	Frequency []domain.Weekday
}

type CreateHabit struct {
	pillarRepo domain.PillarRepository
	habitRepo  domain.HabitRepository
}

func NewCreateHabit(pillarRepo domain.PillarRepository, habitRepo domain.HabitRepository) *CreateHabit {
	return &CreateHabit{pillarRepo: pillarRepo, habitRepo: habitRepo}
}

func (uc *CreateHabit) Execute(input CreateHabitInput) (domain.Habit, error) {
	_, found, err := uc.pillarRepo.FindByID(input.PillarID)
	if err != nil {
		return domain.Habit{}, fmt.Errorf("buscar pilar: %w", err)
	}
	if !found {
		return domain.Habit{}, fmt.Errorf("pilar %d não encontrado", input.PillarID)
	}

	_, exists, err := uc.habitRepo.FindByNameAndPillar(input.Name, input.PillarID)
	if err != nil {
		return domain.Habit{}, fmt.Errorf("verificar duplicidade: %w", err)
	}
	if exists {
		return domain.Habit{}, fmt.Errorf("já existe um hábito com o nome %q neste pilar", input.Name)
	}

	habit, err := domain.NewHabit(input.Name, input.PillarID, input.Color, input.Frequency)
	if err != nil {
		return domain.Habit{}, fmt.Errorf("criar hábito: %w", err)
	}

	saved, err := uc.habitRepo.Save(habit)
	if err != nil {
		return domain.Habit{}, fmt.Errorf("salvar hábito: %w", err)
	}

	return saved, nil
}
