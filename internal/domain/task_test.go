package domain_test

import (
	"testing"
	"time"

	"habits/internal/domain"
)

func TestNewTask_Valid(t *testing.T) {
	task, err := domain.NewTask("Revisar PR", "")
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if task.Title != "Revisar PR" {
		t.Errorf("esperava title %q, got %q", "Revisar PR", task.Title)
	}
	if task.Status != domain.TaskStatusEmAndamento {
		t.Errorf("esperava status em_andamento, got %q", task.Status)
	}
	if task.CreatedAt.IsZero() {
		t.Error("esperava createdAt preenchido")
	}
}

func TestNewTask_TituloVazio(t *testing.T) {
	_, err := domain.NewTask("", "")
	if err == nil {
		t.Fatal("esperava erro para título vazio")
	}
}

func TestTask_TimeSinceUpdate_SemUpdate(t *testing.T) {
	task, _ := domain.NewTask("Tarefa", "")
	task.CreatedAt = time.Now().Add(-2 * time.Hour)

	d := task.TimeSinceUpdate(time.Now())
	if d < time.Hour {
		t.Errorf("esperava duração >= 1h, got %v", d)
	}
}

func TestTask_TimeSinceUpdate_ComUpdate(t *testing.T) {
	task, _ := domain.NewTask("Tarefa", "")
	task.CreatedAt = time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now().Add(-30 * time.Minute)
	task.UpdatedAt = &updatedAt

	d := task.TimeSinceUpdate(time.Now())
	if d >= time.Hour {
		t.Errorf("esperava duração < 1h (usa updatedAt), got %v", d)
	}
}
