package usecase_test

import (
	"testing"

	"habits/internal/domain"
	"habits/internal/usecase"
)

// --- mocks ---

type fakePillarRepo struct {
	pillars map[int64]domain.Pillar
}

func (r *fakePillarRepo) Save(p domain.Pillar) (domain.Pillar, error) {
	p.ID = int64(len(r.pillars) + 1)
	r.pillars[p.ID] = p
	return p, nil
}
func (r *fakePillarRepo) FindByID(id int64) (domain.Pillar, bool, error) {
	p, ok := r.pillars[id]
	return p, ok, nil
}
func (r *fakePillarRepo) FindByName(name string) (domain.Pillar, bool, error) {
	for _, p := range r.pillars {
		if p.Name == name {
			return p, true, nil
		}
	}
	return domain.Pillar{}, false, nil
}
func (r *fakePillarRepo) FindAll() ([]domain.Pillar, error) {
	pillars := make([]domain.Pillar, 0, len(r.pillars))
	for _, p := range r.pillars {
		pillars = append(pillars, p)
	}
	return pillars, nil
}
func (r *fakePillarRepo) Delete(id int64) error              { return nil }

type fakeHabitRepo struct {
	habits map[int64]domain.Habit
}

func (r *fakeHabitRepo) Save(h domain.Habit) (domain.Habit, error) {
	h.ID = int64(len(r.habits) + 1)
	r.habits[h.ID] = h
	return h, nil
}
func (r *fakeHabitRepo) FindByNameAndPillar(name string, pillarID int64) (domain.Habit, bool, error) {
	for _, h := range r.habits {
		if h.Name == name && h.PillarID == pillarID {
			return h, true, nil
		}
	}
	return domain.Habit{}, false, nil
}
func (r *fakeHabitRepo) FindAll() ([]domain.Habit, error) {
	habits := make([]domain.Habit, 0, len(r.habits))
	for _, h := range r.habits {
		habits = append(habits, h)
	}
	return habits, nil
}
func (r *fakeHabitRepo) FindByPillar(pillarID int64) ([]domain.Habit, error) {
	var habits []domain.Habit
	for _, h := range r.habits {
		if h.PillarID == pillarID {
			habits = append(habits, h)
		}
	}
	return habits, nil
}
func (r *fakeHabitRepo) Delete(id int64) error                              { return nil }

func newRepos() (*fakePillarRepo, *fakeHabitRepo) {
	return &fakePillarRepo{pillars: map[int64]domain.Pillar{1: {ID: 1, Name: "Saúde"}}},
		&fakeHabitRepo{habits: map[int64]domain.Habit{}}
}

// --- testes ---

func TestCreateHabit_PilarNaoEncontrado(t *testing.T) {
	pillarRepo, habitRepo := newRepos()
	uc := usecase.NewCreateHabit(pillarRepo, habitRepo)

	_, err := uc.Execute(usecase.CreateHabitInput{
		Name:      "Meditação",
		PillarID:  99,
		Color:     domain.ColorVerde,
		Frequency: []domain.Weekday{domain.WeekdaySeg},
	})
	if err == nil {
		t.Fatal("esperava erro para pilar inexistente")
	}
}

func TestCreateHabit_NomeDuplicado(t *testing.T) {
	pillarRepo, habitRepo := newRepos()
	uc := usecase.NewCreateHabit(pillarRepo, habitRepo)

	input := usecase.CreateHabitInput{
		Name:      "Meditação",
		PillarID:  1,
		Color:     domain.ColorVerde,
		Frequency: []domain.Weekday{domain.WeekdaySeg},
	}
	if _, err := uc.Execute(input); err != nil {
		t.Fatalf("primeira criação falhou: %v", err)
	}
	_, err := uc.Execute(input)
	if err == nil {
		t.Fatal("esperava erro para nome duplicado no mesmo pilar")
	}
}

func TestCreateHabit_Sucesso(t *testing.T) {
	pillarRepo, habitRepo := newRepos()

	uc := usecase.NewCreateHabit(pillarRepo, habitRepo)
	h, err := uc.Execute(usecase.CreateHabitInput{
		Name:      "Meditação",
		PillarID:  1,
		Color:     domain.ColorVerde,
		Frequency: []domain.Weekday{domain.WeekdaySeg, domain.WeekdayQua},
	})

	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if h.ID == 0 {
		t.Error("esperava ID preenchido após salvar")
	}
	if h.Name != "Meditação" {
		t.Errorf("esperava name %q, got %q", "Meditação", h.Name)
	}
}
