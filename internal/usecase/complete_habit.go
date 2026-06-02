package usecase

import (
	"fmt"
	"time"

	"habits/internal/domain"
)

type CompleteHabit struct {
	habitRepo   domain.HabitRepository
	sessionRepo domain.SessionRepository
}

func NewCompleteHabit(habitRepo domain.HabitRepository, sessionRepo domain.SessionRepository) *CompleteHabit {
	return &CompleteHabit{habitRepo: habitRepo, sessionRepo: sessionRepo}
}

func (uc *CompleteHabit) Execute(habitID int64) (domain.Session, error) {
	_, found, err := uc.habitRepo.FindByID(habitID)
	if err != nil {
		return domain.Session{}, fmt.Errorf("buscar hábito: %w", err)
	}
	if !found {
		return domain.Session{}, fmt.Errorf("hábito %d não encontrado", habitID)
	}

	now := time.Now()
	s, err := domain.NewSession(domain.SessionRefHabit, habitID)
	if err != nil {
		return domain.Session{}, err
	}
	s.StartedAt = now
	if err := s.Finish(now); err != nil {
		return domain.Session{}, err
	}

	return uc.sessionRepo.Save(s)
}
