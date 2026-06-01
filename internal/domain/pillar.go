package domain

import (
	"fmt"
	"time"
)

type Pillar struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

func NewPillar(name string) (Pillar, error) {
	if name == "" {
		return Pillar{}, fmt.Errorf("nome não pode ser vazio")
	}
	return Pillar{
		Name:      name,
		CreatedAt: time.Now(),
	}, nil
}
