package domain

import (
	"fmt"
	"time"
)

type SessionRefType string

const (
	SessionRefHabit SessionRefType = "habit"
	SessionRefTask  SessionRefType = "task"
)

type Session struct {
	ID         int64
	RefType    SessionRefType
	RefID      int64
	StartedAt  time.Time
	FinishedAt *time.Time
}

func NewSession(refType SessionRefType, refID int64) (Session, error) {
	if refID == 0 {
		return Session{}, fmt.Errorf("refID não pode ser zero")
	}
	return Session{
		RefType:   refType,
		RefID:     refID,
		StartedAt: time.Now(),
	}, nil
}

func (s *Session) Finish(at time.Time) error {
	if !at.After(s.StartedAt) {
		return fmt.Errorf("finishedAt deve ser posterior ao startedAt")
	}
	s.FinishedAt = &at
	return nil
}
