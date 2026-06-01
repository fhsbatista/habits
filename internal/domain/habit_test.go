package domain_test

import (
	"testing"

	"habits/internal/domain"
)

func TestNewHabit_NomeVazio(t *testing.T) {
	_, err := domain.NewHabit("", 1, domain.ColorAzul, []domain.Weekday{domain.WeekdaySeg})
	if err == nil {
		t.Fatal("esperava erro para nome vazio")
	}
}

func TestNewHabit_FrequenciaVazia(t *testing.T) {
	_, err := domain.NewHabit("Leitura", 1, domain.ColorAzul, []domain.Weekday{})
	if err == nil {
		t.Fatal("esperava erro para frequência vazia")
	}
}

func TestNewHabit_CorInvalida(t *testing.T) {
	_, err := domain.NewHabit("Leitura", 1, domain.Color("roxo"), []domain.Weekday{domain.WeekdaySeg})
	if err == nil {
		t.Fatal("esperava erro para cor inválida")
	}
}

func TestNewHabit_CorAleatoria(t *testing.T) {
	h, err := domain.NewHabit("Leitura", 1, domain.ColorAleatorio, []domain.Weekday{domain.WeekdaySeg})
	if err != nil {
		t.Fatalf("esperava sucesso com cor aleatória, got erro: %v", err)
	}
	validColors := []domain.Color{domain.ColorAzul, domain.ColorVerde, domain.ColorLaranja, domain.ColorVermelho, domain.ColorAmarelo}
	found := false
	for _, c := range validColors {
		if h.Color == c {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("cor aleatória %q não é uma cor válida", h.Color)
	}
}

func TestNewHabit_Valid(t *testing.T) {
	h, err := domain.NewHabit("Meditação", 1, domain.ColorVerde, []domain.Weekday{domain.WeekdaySeg, domain.WeekdayQua})
	if err != nil {
		t.Fatalf("esperava sucesso, got erro: %v", err)
	}
	if h.Name != "Meditação" {
		t.Errorf("esperava name %q, got %q", "Meditação", h.Name)
	}
	if h.PillarID != 1 {
		t.Errorf("esperava pillarID %d, got %d", 1, h.PillarID)
	}
	if h.Color != domain.ColorVerde {
		t.Errorf("esperava color %q, got %q", domain.ColorVerde, h.Color)
	}
	if h.CreatedAt.IsZero() {
		t.Error("esperava createdAt preenchido")
	}
}
