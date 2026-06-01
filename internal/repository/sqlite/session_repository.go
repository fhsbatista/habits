package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"habits/internal/domain"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Save(s domain.Session) (domain.Session, error) {
	res, err := r.db.Exec(
		`INSERT INTO sessions (ref_type, ref_id, started_at, finished_at) VALUES (?, ?, ?, ?)`,
		string(s.RefType), s.RefID,
		s.StartedAt.UTC().Format(time.RFC3339),
		encodeTime(s.FinishedAt),
	)
	if err != nil {
		return domain.Session{}, fmt.Errorf("inserir sessão: %w", err)
	}
	id, _ := res.LastInsertId()
	s.ID = id
	return s, nil
}

func (r *SessionRepository) FindOpen() (domain.Session, bool, error) {
	row := r.db.QueryRow(`SELECT id, ref_type, ref_id, started_at FROM sessions WHERE finished_at IS NULL LIMIT 1`)
	var s domain.Session
	var refType, startedAt string
	err := row.Scan(&s.ID, &refType, &s.RefID, &startedAt)
	if err == sql.ErrNoRows {
		return domain.Session{}, false, nil
	}
	if err != nil {
		return domain.Session{}, false, fmt.Errorf("escanear sessão: %w", err)
	}
	s.RefType = domain.SessionRefType(refType)
	s.StartedAt, _ = time.Parse(time.RFC3339, startedAt)
	return s, true, nil
}

func (r *SessionRepository) FindByDayAndRef(day time.Time, refType domain.SessionRefType, refID int64) ([]domain.Session, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)
	rows, err := r.db.Query(
		`SELECT id, ref_type, ref_id, started_at, finished_at FROM sessions WHERE ref_type=? AND ref_id=? AND started_at >= ? AND started_at < ?`,
		string(refType), refID,
		start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("buscar sessões: %w", err)
	}
	defer rows.Close()
	return scanSessions(rows)
}

func (r *SessionRepository) FindByDay(day time.Time) ([]domain.Session, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	end := start.Add(24 * time.Hour)
	rows, err := r.db.Query(
		`SELECT id, ref_type, ref_id, started_at, finished_at FROM sessions WHERE started_at >= ? AND started_at < ? ORDER BY started_at`,
		start.UTC().Format(time.RFC3339), end.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("buscar sessões do dia: %w", err)
	}
	defer rows.Close()
	return scanSessions(rows)
}

func (r *SessionRepository) Update(s domain.Session) error {
	_, err := r.db.Exec(
		`UPDATE sessions SET finished_at=? WHERE id=?`,
		encodeTime(s.FinishedAt), s.ID,
	)
	return err
}

func scanSessions(rows *sql.Rows) ([]domain.Session, error) {
	var sessions []domain.Session
	for rows.Next() {
		var s domain.Session
		var refType, startedAt string
		var finishedAt sql.NullString
		if err := rows.Scan(&s.ID, &refType, &s.RefID, &startedAt, &finishedAt); err != nil {
			return nil, fmt.Errorf("escanear sessão: %w", err)
		}
		s.RefType = domain.SessionRefType(refType)
		s.StartedAt, _ = time.Parse(time.RFC3339, startedAt)
		if finishedAt.Valid && finishedAt.String != "" {
			t, _ := time.Parse(time.RFC3339, finishedAt.String)
			s.FinishedAt = &t
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}
