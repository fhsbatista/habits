package usecase_test

import (
	"testing"

	"habits/internal/domain"
	"habits/internal/usecase"
)

func TestRemoveHabit_Sucesso(t *testing.T) {
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{
		1: {ID: 1, Name: "Meditação", PillarID: 1},
	}}
	uc := usecase.NewRemoveHabit(habitRepo)

	if err := uc.Execute(1); err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if _, ok := habitRepo.habits[1]; ok {
		t.Error("esperava hábito removido")
	}
}

func TestRemoveHabit_NaoEncontrado(t *testing.T) {
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{}}
	uc := usecase.NewRemoveHabit(habitRepo)

	if err := uc.Execute(99); err == nil {
		t.Fatal("esperava erro para hábito inexistente")
	}
}
