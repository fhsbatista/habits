package domain_test

import (
	"testing"

	"habits/internal/domain"
)

func TestNewPillar_NomeVazio(t *testing.T) {
	_, err := domain.NewPillar("")
	if err == nil {
		t.Fatal("esperava erro para nome vazio")
	}
}

func TestNewPillar_Valid(t *testing.T) {
	p, err := domain.NewPillar("Saúde")
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if p.Name != "Saúde" {
		t.Errorf("esperava name %q, got %q", "Saúde", p.Name)
	}
	if p.CreatedAt.IsZero() {
		t.Error("esperava createdAt preenchido")
	}
}
