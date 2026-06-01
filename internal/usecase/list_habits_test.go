package usecase_test

import (
	"testing"

	"habits/internal/domain"
	"habits/internal/usecase"
)

func TestListHabits_RetornaTodos(t *testing.T) {
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{
		1: {ID: 1, Name: "Meditação", PillarID: 1},
		2: {ID: 2, Name: "Leitura", PillarID: 2},
	}}
	uc := usecase.NewListHabits(habitRepo)

	habits, err := uc.Execute(0)
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if len(habits) != 2 {
		t.Errorf("esperava 2 hábitos, got %d", len(habits))
	}
}

func TestListHabits_PorPilar(t *testing.T) {
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{
		1: {ID: 1, Name: "Meditação", PillarID: 1},
		2: {ID: 2, Name: "Leitura", PillarID: 2},
	}}
	uc := usecase.NewListHabits(habitRepo)

	habits, err := uc.Execute(1)
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if len(habits) != 1 {
		t.Errorf("esperava 1 hábito, got %d", len(habits))
	}
	if habits[0].PillarID != 1 {
		t.Errorf("esperava pillarID 1, got %d", habits[0].PillarID)
	}
}
