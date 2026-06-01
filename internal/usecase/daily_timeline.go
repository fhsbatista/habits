package usecase

import (
	"fmt"
	"time"

	"habits/internal/domain"
)

type TimelineEntry struct {
	Session      domain.Session
	Name         string
	RefType      domain.SessionRefType
	Color        domain.Color
	EffectiveEnd time.Time
	IsOpen       bool
}

type DailyTimeline struct {
	habitRepo   domain.HabitRepository
	taskRepo    domain.TaskRepository
	sessionRepo domain.SessionRepository
}

func NewDailyTimeline(habitRepo domain.HabitRepository, taskRepo domain.TaskRepository, sessionRepo domain.SessionRepository) *DailyTimeline {
	return &DailyTimeline{habitRepo: habitRepo, taskRepo: taskRepo, sessionRepo: sessionRepo}
}

func (uc *DailyTimeline) Execute(day time.Time) ([]TimelineEntry, error) {
	sessions, err := uc.sessionRepo.FindByDay(day)
	if err != nil {
		return nil, fmt.Errorf("buscar sessões: %w", err)
	}

	entries := make([]TimelineEntry, 0, len(sessions))
	for _, s := range sessions {
		entry, err := uc.buildEntry(s, day)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (uc *DailyTimeline) buildEntry(s domain.Session, now time.Time) (TimelineEntry, error) {
	entry := TimelineEntry{
		Session: s,
		RefType: s.RefType,
	}

	if s.FinishedAt != nil {
		entry.EffectiveEnd = *s.FinishedAt
	} else {
		entry.EffectiveEnd = now
		entry.IsOpen = true
	}

	switch s.RefType {
	case domain.SessionRefHabit:
		h, found, err := uc.habitRepo.FindByID(s.RefID)
		if err != nil {
			return TimelineEntry{}, fmt.Errorf("buscar hábito %d: %w", s.RefID, err)
		}
		if found {
			entry.Name = h.Name
			entry.Color = h.Color
		} else {
			entry.Name = fmt.Sprintf("hábito %d", s.RefID)
		}
	case domain.SessionRefTask:
		t, found, err := uc.taskRepo.FindByID(s.RefID)
		if err != nil {
			return TimelineEntry{}, fmt.Errorf("buscar tarefa %d: %w", s.RefID, err)
		}
		if found {
			entry.Name = t.Title
		} else {
			entry.Name = fmt.Sprintf("tarefa %d", s.RefID)
		}
	}

	return entry, nil
}
