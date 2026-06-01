package domain

import (
	"fmt"
	"time"
)

type TaskStatus string

const (
	TaskStatusEmAndamento TaskStatus = "em_andamento"
	TaskStatusConcluida   TaskStatus = "concluida"
	TaskStatusArquivada   TaskStatus = "arquivada"
)

type NextActionType string

const (
	NextActionExecutar NextActionType = "executar"
	NextActionAguardar NextActionType = "aguardar"
)

type NextAction struct {
	Type NextActionType
}

type Task struct {
	ID          int64
	Title       string
	Description string
	Status      TaskStatus
	NextAction  *NextAction
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

func NewTask(title, description string) (Task, error) {
	if title == "" {
		return Task{}, fmt.Errorf("título não pode ser vazio")
	}
	return Task{
		Title:       title,
		Description: description,
		Status:      TaskStatusEmAndamento,
		CreatedAt:   time.Now(),
	}, nil
}

func (t Task) TimeSinceUpdate(now time.Time) time.Duration {
	if t.UpdatedAt != nil {
		return now.Sub(*t.UpdatedAt)
	}
	return now.Sub(t.CreatedAt)
}
