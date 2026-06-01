package usecase_test

import (
	"testing"
	"time"

	"habits/internal/domain"
	"habits/internal/usecase"
)

type fakeSessionRepo struct {
	sessions map[int64]domain.Session
}

func (r *fakeSessionRepo) Save(s domain.Session) (domain.Session, error) {
	s.ID = int64(len(r.sessions) + 1)
	r.sessions[s.ID] = s
	return s, nil
}
func (r *fakeSessionRepo) FindOpen() (domain.Session, bool, error) {
	for _, s := range r.sessions {
		if s.FinishedAt == nil {
			return s, true, nil
		}
	}
	return domain.Session{}, false, nil
}
func (r *fakeSessionRepo) FindByDayAndRef(day time.Time, refType domain.SessionRefType, refID int64) ([]domain.Session, error) {
	return nil, nil
}
func (r *fakeSessionRepo) FindByDay(day time.Time) ([]domain.Session, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)
	var result []domain.Session
	for _, s := range r.sessions {
		if s.StartedAt.After(start) && s.StartedAt.Before(end) {
			result = append(result, s)
		}
	}
	return result, nil
}
func (r *fakeSessionRepo) FindByDateRange(start, end time.Time) ([]domain.Session, error) {
	var result []domain.Session
	for _, s := range r.sessions {
		if (s.StartedAt.Equal(start) || s.StartedAt.After(start)) && s.StartedAt.Before(end) {
			result = append(result, s)
		}
	}
	return result, nil
}
func (r *fakeSessionRepo) Update(s domain.Session) error {
	r.sessions[s.ID] = s
	return nil
}

func TestCheckIn_Sucesso(t *testing.T) {
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{
		1: {ID: 1, Name: "Meditação"},
	}}
	sessionRepo := &fakeSessionRepo{sessions: map[int64]domain.Session{}}
	uc := usecase.NewCheckIn(habitRepo, nil, sessionRepo)

	s, err := uc.Execute(domain.SessionRefHabit, 1)
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if s.ID == 0 {
		t.Error("esperava ID preenchido após salvar")
	}
	if s.FinishedAt != nil {
		t.Error("esperava sessão aberta")
	}
}

func TestCheckIn_SessaoJaAberta(t *testing.T) {
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{
		1: {ID: 1, Name: "Meditação"},
	}}
	openSession := domain.Session{ID: 1, RefType: domain.SessionRefHabit, RefID: 1, StartedAt: time.Now()}
	sessionRepo := &fakeSessionRepo{sessions: map[int64]domain.Session{1: openSession}}
	uc := usecase.NewCheckIn(habitRepo, nil, sessionRepo)

	_, err := uc.Execute(domain.SessionRefHabit, 1)
	if err == nil {
		t.Fatal("esperava erro para sessão já aberta")
	}
}

func TestCheckIn_HabitoNaoEncontrado(t *testing.T) {
	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{}}
	sessionRepo := &fakeSessionRepo{sessions: map[int64]domain.Session{}}
	uc := usecase.NewCheckIn(habitRepo, nil, sessionRepo)

	_, err := uc.Execute(domain.SessionRefHabit, 99)
	if err == nil {
		t.Fatal("esperava erro para hábito inexistente")
	}
}
