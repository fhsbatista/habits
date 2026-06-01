package usecase_test

import (
	"testing"
	"time"

	"habits/internal/domain"
	"habits/internal/usecase"
)

type fakeTaskRepo struct {
	tasks map[int64]domain.Task
}

func (r *fakeTaskRepo) Save(t domain.Task) (domain.Task, error) {
	t.ID = int64(len(r.tasks) + 1)
	r.tasks[t.ID] = t
	return t, nil
}
func (r *fakeTaskRepo) FindByTitle(title string) (domain.Task, bool, error) {
	for _, t := range r.tasks {
		if t.Title == title {
			return t, true, nil
		}
	}
	return domain.Task{}, false, nil
}
func (r *fakeTaskRepo) FindInProgress() ([]domain.Task, error) {
	var tasks []domain.Task
	for _, t := range r.tasks {
		if t.Status == domain.TaskStatusEmAndamento {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}
func (r *fakeTaskRepo) Update(t domain.Task) error { r.tasks[t.ID] = t; return nil }
func (r *fakeTaskRepo) Delete(id int64) error      { delete(r.tasks, id); return nil }

func TestDailyOverview_SeparaTarefasPorNexAction(t *testing.T) {
	now := time.Now()
	executar := domain.NextActionType("executar")
	aguardar := domain.NextActionType("aguardar")

	taskRepo := &fakeTaskRepo{tasks: map[int64]domain.Task{
		1: {ID: 1, Title: "Tarefa A", Status: domain.TaskStatusEmAndamento, NextAction: &domain.NextAction{Type: executar}, CreatedAt: now},
		2: {ID: 2, Title: "Tarefa B", Status: domain.TaskStatusEmAndamento, NextAction: &domain.NextAction{Type: aguardar}, CreatedAt: now},
		3: {ID: 3, Title: "Tarefa C", Status: domain.TaskStatusEmAndamento, NextAction: nil, CreatedAt: now},
		4: {ID: 4, Title: "Tarefa D", Status: domain.TaskStatusConcluida, NextAction: nil, CreatedAt: now},
	}}

	uc := usecase.NewDailyOverview(
		&fakePillarRepo{pillars: map[int64]domain.Pillar{}},
		&fakeHabitRepo{habits: map[int64]domain.Habit{}},
		taskRepo,
	)

	result, err := uc.Execute(now)
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if len(result.ExecuteTasks) != 2 {
		t.Errorf("esperava 2 tarefas em executar (inclui sem next action), got %d", len(result.ExecuteTasks))
	}
	if len(result.WaitTasks) != 1 {
		t.Errorf("esperava 1 tarefa em aguardar, got %d", len(result.WaitTasks))
	}
}

func TestDailyOverview_OrdenaPorTempoSemUpdate(t *testing.T) {
	now := time.Now()
	executar := domain.NextActionType("executar")
	maisAntiga := now.Add(-5 * 24 * time.Hour)
	maisRecente := now.Add(-1 * time.Hour)

	taskRepo := &fakeTaskRepo{tasks: map[int64]domain.Task{
		1: {ID: 1, Title: "Recente", Status: domain.TaskStatusEmAndamento, NextAction: &domain.NextAction{Type: executar}, CreatedAt: maisRecente},
		2: {ID: 2, Title: "Antiga", Status: domain.TaskStatusEmAndamento, NextAction: &domain.NextAction{Type: executar}, CreatedAt: maisAntiga},
	}}

	uc := usecase.NewDailyOverview(
		&fakePillarRepo{pillars: map[int64]domain.Pillar{}},
		&fakeHabitRepo{habits: map[int64]domain.Habit{}},
		taskRepo,
	)

	result, err := uc.Execute(now)
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if result.ExecuteTasks[0].Title != "Antiga" {
		t.Errorf("esperava 'Antiga' primeiro, got %q", result.ExecuteTasks[0].Title)
	}
}

func TestDailyOverview_HabitosPendentesHoje(t *testing.T) {
	now := time.Now()
	today := domain.WeekdayFromTime(now)

	habitRepo := &fakeHabitRepo{habits: map[int64]domain.Habit{
		1: {ID: 1, Name: "Meditação", PillarID: 1, Color: domain.ColorVerde, Frequency: []domain.Weekday{today}},
	}}
	pillarRepo := &fakePillarRepo{pillars: map[int64]domain.Pillar{
		1: {ID: 1, Name: "Saúde"},
	}}

	uc := usecase.NewDailyOverview(pillarRepo, habitRepo, &fakeTaskRepo{tasks: map[int64]domain.Task{}})

	result, err := uc.Execute(now)
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if len(result.PendingHabits) != 1 {
		t.Errorf("esperava 1 hábito pendente, got %d", len(result.PendingHabits))
	}
	if result.PendingHabits[0].Habit.Name != "Meditação" {
		t.Errorf("esperava hábito 'Meditação', got %q", result.PendingHabits[0].Habit.Name)
	}
}
