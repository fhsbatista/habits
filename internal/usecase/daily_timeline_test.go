package usecase_test

import (
	"testing"
	"time"

	"habits/internal/domain"
	"habits/internal/usecase"
)

func TestDailyTimeline_RetornaEntradas(t *testing.T) {
	now := time.Now()
	started := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	finished := time.Date(now.Year(), now.Month(), now.Day(), 10, 30, 0, 0, now.Location())

	sessionRepo := &fakeSessionRepo{sessions: map[int64]domain.Session{
		1: {ID: 1, RefType: domain.SessionRefHabit, RefID: 1, StartedAt: started, FinishedAt: &finished},
	}}
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{
		1: {ID: 1, Name: "Meditação", Color: domain.ColorVerde},
	}}

	uc := usecase.NewDailyTimeline(habitRepo, &fakeTaskRepo{tasks: map[int64]domain.Task{}}, sessionRepo)
	entries, err := uc.Execute(now)
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("esperava 1 entrada, got %d", len(entries))
	}
	if entries[0].Name != "Meditação" {
		t.Errorf("esperava name %q, got %q", "Meditação", entries[0].Name)
	}
	if entries[0].Color != domain.ColorVerde {
		t.Errorf("esperava cor verde, got %q", entries[0].Color)
	}
}

func TestDailyTimeline_SessaoAbertaUsaNow(t *testing.T) {
	now := time.Now()
	started := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())

	sessionRepo := &fakeSessionRepo{sessions: map[int64]domain.Session{
		1: {ID: 1, RefType: domain.SessionRefHabit, RefID: 1, StartedAt: started, FinishedAt: nil},
	}}
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{
		1: {ID: 1, Name: "Exercício", Color: domain.ColorAzul},
	}}

	uc := usecase.NewDailyTimeline(habitRepo, &fakeTaskRepo{tasks: map[int64]domain.Task{}}, sessionRepo)
	entries, err := uc.Execute(now)
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if entries[0].EffectiveEnd.IsZero() {
		t.Error("esperava EffectiveEnd preenchido para sessão aberta")
	}
	if entries[0].IsOpen != true {
		t.Error("esperava IsOpen true para sessão sem finishedAt")
	}
}
