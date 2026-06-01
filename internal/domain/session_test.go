package domain_test

import (
	"testing"
	"time"

	"habits/internal/domain"
)

func TestNewSession_Valid(t *testing.T) {
	s, err := domain.NewSession(domain.SessionRefHabit, 1)
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if s.RefType != domain.SessionRefHabit {
		t.Errorf("esperava refType %q, got %q", domain.SessionRefHabit, s.RefType)
	}
	if s.RefID != 1 {
		t.Errorf("esperava refID 1, got %d", s.RefID)
	}
	if s.StartedAt.IsZero() {
		t.Error("esperava startedAt preenchido")
	}
	if s.FinishedAt != nil {
		t.Error("esperava finishedAt nulo ao criar")
	}
}

func TestNewSession_RefIDZero(t *testing.T) {
	_, err := domain.NewSession(domain.SessionRefHabit, 0)
	if err == nil {
		t.Fatal("esperava erro para refID zero")
	}
}

func TestSession_Finish(t *testing.T) {
	s, _ := domain.NewSession(domain.SessionRefHabit, 1)
	now := time.Now().Add(time.Minute)

	if err := s.Finish(now); err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if s.FinishedAt == nil {
		t.Error("esperava finishedAt preenchido")
	}
}

func TestSession_Finish_AntesDoinicio(t *testing.T) {
	s, _ := domain.NewSession(domain.SessionRefHabit, 1)
	before := s.StartedAt.Add(-time.Minute)

	if err := s.Finish(before); err == nil {
		t.Fatal("esperava erro para finishedAt anterior ao startedAt")
	}
}
