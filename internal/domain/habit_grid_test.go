package domain_test

import (
	"testing"
	"time"

	"habits/internal/domain"
)

func TestHabit_WasPerformedOn_Sim(t *testing.T) {
	day := time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)
	started := time.Date(2026, 6, 1, 10, 0, 0, 0, time.Local)
	finished := time.Date(2026, 6, 1, 10, 30, 0, 0, time.Local)

	sessions := []domain.Session{
		{ID: 1, RefType: domain.SessionRefHabit, RefID: 1, StartedAt: started, FinishedAt: &finished},
	}

	h, _ := domain.NewHabit("Meditação", 1, domain.ColorVerde, []domain.Weekday{domain.WeekdaySeg})
	h.ID = 1

	if !h.WasPerformedOn(day, sessions) {
		t.Error("esperava que hábito tivesse sido performado no dia")
	}
}

func TestHabit_WasPerformedOn_Nao(t *testing.T) {
	day := time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)
	h, _ := domain.NewHabit("Meditação", 1, domain.ColorVerde, []domain.Weekday{domain.WeekdaySeg})
	h.ID = 1

	if h.WasPerformedOn(day, []domain.Session{}) {
		t.Error("esperava que hábito NÃO tivesse sido performado no dia")
	}
}

func TestHabit_WasPerformedOn_SessaoAbertaNaoContabiliza(t *testing.T) {
	day := time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)
	started := time.Date(2026, 6, 1, 10, 0, 0, 0, time.Local)

	sessions := []domain.Session{
		{ID: 1, RefType: domain.SessionRefHabit, RefID: 1, StartedAt: started, FinishedAt: nil},
	}

	h, _ := domain.NewHabit("Meditação", 1, domain.ColorVerde, []domain.Weekday{domain.WeekdaySeg})
	h.ID = 1

	if h.WasPerformedOn(day, sessions) {
		t.Error("sessão aberta não deve contabilizar como performado")
	}
}
