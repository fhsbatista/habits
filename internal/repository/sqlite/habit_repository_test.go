package sqlite_test

import (
	"testing"

	"habits/internal/domain"
	"habits/internal/repository/sqlite"
)

func openTestDB(t *testing.T) *sqlite.HabitRepository {
	t.Helper()
	db, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("abrir banco: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return sqlite.NewHabitRepository(db)
}

func TestHabitRepository_FindByNameAndPillar(t *testing.T) {
	repo := openTestDB(t)

	habit, _ := domain.NewHabit("Meditação", 1, domain.ColorVerde, []domain.Weekday{domain.WeekdaySeg})
	repo.Save(habit)

	found, ok, err := repo.FindByNameAndPillar("Meditação", 1)
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if !ok {
		t.Fatal("esperava encontrar o hábito")
	}
	if found.Name != "Meditação" {
		t.Errorf("esperava name %q, got %q", "Meditação", found.Name)
	}

	_, ok, _ = repo.FindByNameAndPillar("Meditação", 2)
	if ok {
		t.Error("não devia encontrar hábito com pillar diferente")
	}
}

func TestHabitRepository_Save(t *testing.T) {
	repo := openTestDB(t)

	habit, _ := domain.NewHabit("Meditação", 1, domain.ColorVerde, []domain.Weekday{domain.WeekdaySeg})
	saved, err := repo.Save(habit)
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if saved.ID == 0 {
		t.Error("esperava ID preenchido após salvar")
	}
}
