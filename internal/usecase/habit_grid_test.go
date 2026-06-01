package usecase_test

import (
	"testing"
	"time"

	"habits/internal/domain"
	"habits/internal/usecase"
)

func TestHabitGrid_RetornaGridComDiasDoMes(t *testing.T) {
	h1 := domain.Habit{ID: 1, Name: "Meditação", Color: domain.ColorVerde}
	h2 := domain.Habit{ID: 2, Name: "Exercício", Color: domain.ColorAzul}

	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{1: h1, 2: h2}}

	day1 := time.Date(2026, 6, 1, 10, 0, 0, 0, time.Local)
	fin1 := day1.Add(30 * time.Minute)
	day3 := time.Date(2026, 6, 3, 9, 0, 0, 0, time.Local)
	fin3 := day3.Add(20 * time.Minute)

	sessionRepo := &fakeSessionRepo{sessions: map[int64]domain.Session{
		1: {ID: 1, RefType: domain.SessionRefHabit, RefID: 1, StartedAt: day1, FinishedAt: &fin1},
		2: {ID: 2, RefType: domain.SessionRefHabit, RefID: 2, StartedAt: day3, FinishedAt: &fin3},
	}}

	uc := usecase.NewHabitGrid(habitRepo, sessionRepo)
	month := time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)
	result, err := uc.Execute(month, month.AddDate(0, 1, 0))
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if len(result.Habits) != 2 {
		t.Fatalf("esperava 2 hábitos, got %d", len(result.Habits))
	}
	if len(result.Days) != 30 {
		t.Fatalf("esperava 30 dias em junho, got %d", len(result.Days))
	}

	// habit 1 performed on day 1
	if !result.Performed[1][time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)] {
		t.Error("hábito 1 deveria estar marcado no dia 01/06")
	}
	// habit 1 not performed on day 3
	if result.Performed[1][time.Date(2026, 6, 3, 0, 0, 0, 0, time.Local)] {
		t.Error("hábito 1 não deveria estar marcado no dia 03/06")
	}
	// habit 2 performed on day 3
	if !result.Performed[2][time.Date(2026, 6, 3, 0, 0, 0, 0, time.Local)] {
		t.Error("hábito 2 deveria estar marcado no dia 03/06")
	}
}
