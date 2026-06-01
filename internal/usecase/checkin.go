package usecase

import (
	"fmt"

	"habits/internal/domain"
)

type CheckIn struct {
	habitRepo   domain.HabitRepository
	taskRepo    domain.TaskRepository
	sessionRepo domain.SessionRepository
}

func NewCheckIn(habitRepo domain.HabitRepository, taskRepo domain.TaskRepository, sessionRepo domain.SessionRepository) *CheckIn {
	return &CheckIn{habitRepo: habitRepo, taskRepo: taskRepo, sessionRepo: sessionRepo}
}

func (uc *CheckIn) Execute(refType domain.SessionRefType, refID int64) (domain.Session, error) {
	_, open, err := uc.sessionRepo.FindOpen()
	if err != nil {
		return domain.Session{}, fmt.Errorf("verificar sessão aberta: %w", err)
	}
	if open {
		return domain.Session{}, fmt.Errorf("já existe uma sessão em andamento; execute 'habits stop' antes de iniciar outra")
	}

	if err := uc.validateRef(refType, refID); err != nil {
		return domain.Session{}, err
	}

	session, err := domain.NewSession(refType, refID)
	if err != nil {
		return domain.Session{}, fmt.Errorf("criar sessão: %w", err)
	}

	saved, err := uc.sessionRepo.Save(session)
	if err != nil {
		return domain.Session{}, fmt.Errorf("salvar sessão: %w", err)
	}

	return saved, nil
}

func (uc *CheckIn) validateRef(refType domain.SessionRefType, refID int64) error {
	switch refType {
	case domain.SessionRefHabit:
		_, found, err := uc.habitRepo.FindByID(refID)
		if err != nil {
			return fmt.Errorf("buscar hábito: %w", err)
		}
		if !found {
			return fmt.Errorf("hábito %d não encontrado", refID)
		}
	case domain.SessionRefTask:
		_, found, err := uc.taskRepo.FindByID(refID)
		if err != nil {
			return fmt.Errorf("buscar tarefa: %w", err)
		}
		if !found {
			return fmt.Errorf("tarefa %d não encontrada", refID)
		}
	}
	return nil
}
