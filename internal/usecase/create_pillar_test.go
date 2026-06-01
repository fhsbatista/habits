package usecase_test

import (
	"testing"

	"habits/internal/domain"
	"habits/internal/usecase"
)

func TestCreatePillar_NomeDuplicado(t *testing.T) {
	pillarRepo := &fakePillarRepo{pillars: map[int64]domain.Pillar{}}
	uc := usecase.NewCreatePillar(pillarRepo)

	if _, err := uc.Execute("Saúde"); err != nil {
		t.Fatalf("primeira criação falhou: %v", err)
	}
	_, err := uc.Execute("Saúde")
	if err == nil {
		t.Fatal("esperava erro para nome duplicado")
	}
}

func TestCreatePillar_Sucesso(t *testing.T) {
	pillarRepo := &fakePillarRepo{pillars: map[int64]domain.Pillar{}}
	uc := usecase.NewCreatePillar(pillarRepo)

	p, err := uc.Execute("Saúde")
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if p.ID == 0 {
		t.Error("esperava ID preenchido após salvar")
	}
	if p.Name != "Saúde" {
		t.Errorf("esperava name %q, got %q", "Saúde", p.Name)
	}
}
