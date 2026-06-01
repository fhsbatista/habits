package usecase_test

import (
	"testing"
	"time"

	"habits/internal/domain"
	"habits/internal/usecase"
)

func TestCheckOut_Sucesso(t *testing.T) {
	open := domain.Session{ID: 1, RefType: domain.SessionRefHabit, RefID: 1, StartedAt: time.Now().Add(-10 * time.Minute)}
	sessionRepo := &fakeSessionRepo{sessions: map[int64]domain.Session{1: open}}
	uc := usecase.NewCheckOut(sessionRepo)

	s, err := uc.Execute()
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if s.FinishedAt == nil {
		t.Error("esperava finishedAt preenchido")
	}
}

func TestCheckOut_SemSessaoAberta(t *testing.T) {
	sessionRepo := &fakeSessionRepo{sessions: map[int64]domain.Session{}}
	uc := usecase.NewCheckOut(sessionRepo)

	_, err := uc.Execute()
	if err == nil {
		t.Fatal("esperava erro para nenhuma sessão aberta")
	}
}
