package usecase

import (
	"fmt"

	"habits/internal/domain"
)

type CreatePillar struct {
	pillarRepo domain.PillarRepository
}

func NewCreatePillar(pillarRepo domain.PillarRepository) *CreatePillar {
	return &CreatePillar{pillarRepo: pillarRepo}
}

func (uc *CreatePillar) Execute(name string) (domain.Pillar, error) {
	_, exists, err := uc.pillarRepo.FindByName(name)
	if err != nil {
		return domain.Pillar{}, fmt.Errorf("verificar duplicidade: %w", err)
	}
	if exists {
		return domain.Pillar{}, fmt.Errorf("já existe um pilar com o nome %q", name)
	}

	pillar, err := domain.NewPillar(name)
	if err != nil {
		return domain.Pillar{}, fmt.Errorf("criar pilar: %w", err)
	}

	saved, err := uc.pillarRepo.Save(pillar)
	if err != nil {
		return domain.Pillar{}, fmt.Errorf("salvar pilar: %w", err)
	}

	return saved, nil
}
