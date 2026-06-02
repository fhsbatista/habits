package usecase_test

import (
	"testing"

	"habits/internal/domain"
	"habits/internal/usecase"
)

func TestCompleteHabit_CriaSessionJaFinalizada(t *testing.T) {
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{
		1: {ID: 1, Name: "Meditação"},
	}}
	sessionRepo := &fakeSessionRepo{sessions: map[int64]domain.Session{}}

	uc := usecase.NewCompleteHabit(habitRepo, sessionRepo)
	s, err := uc.Execute(1)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if s.RefType != domain.SessionRefHabit {
		t.Errorf("esperava RefType habit, got %q", s.RefType)
	}
	if s.RefID != 1 {
		t.Errorf("esperava RefID 1, got %d", s.RefID)
	}
	if s.FinishedAt == nil {
		t.Error("sessão deveria já estar finalizada")
	}
}

func TestCompleteHabit_HabitoInexistente(t *testing.T) {
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{}}
	sessionRepo := &fakeSessionRepo{sessions: map[int64]domain.Session{}}

	uc := usecase.NewCompleteHabit(habitRepo, sessionRepo)
	_, err := uc.Execute(99)
	if err == nil {
		t.Error("esperava erro para hábito inexistente")
	}
}
