package domain

import (
	"fmt"
	"math/rand"
	"time"
)

type Color string

const (
	ColorAzul      Color = "azul"
	ColorVerde     Color = "verde"
	ColorLaranja   Color = "laranja"
	ColorVermelho  Color = "vermelho"
	ColorAmarelo   Color = "amarelo"
	ColorAleatorio Color = ""
)

var validColors = []Color{ColorAzul, ColorVerde, ColorLaranja, ColorVermelho, ColorAmarelo}

type Weekday string

const (
	WeekdayDom Weekday = "dom"
	WeekdaySeg Weekday = "seg"
	WeekdayTer Weekday = "ter"
	WeekdayQua Weekday = "qua"
	WeekdayQui Weekday = "qui"
	WeekdaySex Weekday = "sex"
	WeekdaySab Weekday = "sab"
)

type Habit struct {
	ID        int64
	Name      string
	PillarID  int64
	Color     Color
	Frequency []Weekday
	CreatedAt time.Time
}

func NewHabit(name string, pillarID int64, color Color, frequency []Weekday) (Habit, error) {
	if name == "" {
		return Habit{}, fmt.Errorf("nome não pode ser vazio")
	}
	if len(frequency) == 0 {
		return Habit{}, fmt.Errorf("frequência não pode ser vazia")
	}
	if color == ColorAleatorio {
		color = validColors[rand.Intn(len(validColors))]
	} else if !isValidColor(color) {
		return Habit{}, fmt.Errorf("cor inválida: %q", color)
	}
	return Habit{
		Name:      name,
		PillarID:  pillarID,
		Color:     color,
		Frequency: frequency,
		CreatedAt: time.Now(),
	}, nil
}

func isValidColor(c Color) bool {
	for _, v := range validColors {
		if c == v {
			return true
		}
	}
	return false
}
