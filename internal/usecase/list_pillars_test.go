package usecase_test

import (
	"testing"

	"habits/internal/domain"
	"habits/internal/usecase"
)

func TestListPillars_RetornaTodos(t *testing.T) {
	pillarRepo := &fakePillarRepo{pillars: map[int64]domain.Pillar{
		1: {ID: 1, Name: "Saúde"},
		2: {ID: 2, Name: "Carreira"},
	}}
	uc := usecase.NewListPillars(pillarRepo)

	pillars, err := uc.Execute()
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if len(pillars) != 2 {
		t.Errorf("esperava 2 pilares, got %d", len(pillars))
	}
}
