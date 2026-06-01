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

func (h Habit) IsDueOn(day Weekday) bool {
	for _, d := range h.Frequency {
		if d == day {
			return true
		}
	}
	return false
}

func (h Habit) IsDueToday() bool {
	return h.IsDueOn(WeekdayFromTime(time.Now()))
}

func WeekdayFromTime(t time.Time) Weekday {
	switch t.Weekday() {
	case time.Sunday:
		return WeekdayDom
	case time.Monday:
		return WeekdaySeg
	case time.Tuesday:
		return WeekdayTer
	case time.Wednesday:
		return WeekdayQua
	case time.Thursday:
		return WeekdayQui
	case time.Friday:
		return WeekdaySex
	default:
		return WeekdaySab
	}
}

func isValidColor(c Color) bool {
	for _, v := range validColors {
		if c == v {
			return true
		}
	}
	return false
}
