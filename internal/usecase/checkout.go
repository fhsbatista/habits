package usecase

import (
	"fmt"
	"time"

	"habits/internal/domain"
)

type CheckOut struct {
	sessionRepo domain.SessionRepository
}

func NewCheckOut(sessionRepo domain.SessionRepository) *CheckOut {
	return &CheckOut{sessionRepo: sessionRepo}
}

func (uc *CheckOut) Execute() (domain.Session, error) {
	session, found, err := uc.sessionRepo.FindOpen()
	if err != nil {
		return domain.Session{}, fmt.Errorf("buscar sessão aberta: %w", err)
	}
	if !found {
		return domain.Session{}, fmt.Errorf("nenhuma sessão em andamento; use 'habits start' para iniciar uma")
	}

	if err := session.Finish(time.Now()); err != nil {
		return domain.Session{}, fmt.Errorf("finalizar sessão: %w", err)
	}

	if err := uc.sessionRepo.Update(session); err != nil {
		return domain.Session{}, fmt.Errorf("salvar sessão: %w", err)
	}

	return session, nil
}
