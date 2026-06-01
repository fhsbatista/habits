package usecase

import (
	"fmt"

	"habits/internal/domain"
)

type ListPillars struct {
	pillarRepo domain.PillarRepository
}

func NewListPillars(pillarRepo domain.PillarRepository) *ListPillars {
	return &ListPillars{pillarRepo: pillarRepo}
}

func (uc *ListPillars) Execute() ([]domain.Pillar, error) {
	pillars, err := uc.pillarRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("listar pilares: %w", err)
	}
	return pillars, nil
}
