package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"habits/internal/domain"
)

type HabitRepository struct {
	db *sql.DB
}

func NewHabitRepository(db *sql.DB) *HabitRepository {
	return &HabitRepository{db: db}
}

func (r *HabitRepository) Save(h domain.Habit) (domain.Habit, error) {
	freq := encodeFrequency(h.Frequency)
	res, err := r.db.Exec(
		`INSERT INTO habits (name, pillar_id, color, frequency, created_at) VALUES (?, ?, ?, ?, ?)`,
		h.Name, h.PillarID, string(h.Color), freq, h.CreatedAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return domain.Habit{}, fmt.Errorf("inserir hábito: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Habit{}, fmt.Errorf("obter id: %w", err)
	}
	h.ID = id
	return h, nil
}

func (r *HabitRepository) FindByID(id int64) (domain.Habit, bool, error) {
	row := r.db.QueryRow(
		`SELECT id, name, pillar_id, color, frequency, created_at FROM habits WHERE id = ?`, id,
	)
	return scanHabit(row)
}

func (r *HabitRepository) FindByNameAndPillar(name string, pillarID int64) (domain.Habit, bool, error) {
	row := r.db.QueryRow(
		`SELECT id, name, pillar_id, color, frequency, created_at FROM habits WHERE name = ? AND pillar_id = ?`,
		name, pillarID,
	)
	return scanHabit(row)
}

func (r *HabitRepository) FindAll() ([]domain.Habit, error) {
	rows, err := r.db.Query(`SELECT id, name, pillar_id, color, frequency, created_at FROM habits`)
	if err != nil {
		return nil, fmt.Errorf("listar hábitos: %w", err)
	}
	defer rows.Close()
	return scanHabits(rows)
}

func (r *HabitRepository) FindByPillar(pillarID int64) ([]domain.Habit, error) {
	rows, err := r.db.Query(
		`SELECT id, name, pillar_id, color, frequency, created_at FROM habits WHERE pillar_id = ?`, pillarID,
	)
	if err != nil {
		return nil, fmt.Errorf("listar hábitos por pilar: %w", err)
	}
	defer rows.Close()
	return scanHabits(rows)
}

func (r *HabitRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM habits WHERE id = ?`, id)
	return err
}

func scanHabit(row *sql.Row) (domain.Habit, bool, error) {
	var h domain.Habit
	var color, freq, createdAt string
	err := row.Scan(&h.ID, &h.Name, &h.PillarID, &color, &freq, &createdAt)
	if err == sql.ErrNoRows {
		return domain.Habit{}, false, nil
	}
	if err != nil {
		return domain.Habit{}, false, fmt.Errorf("escanear hábito: %w", err)
	}
	h.Color = domain.Color(color)
	h.Frequency = decodeFrequency(freq)
	h.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return h, true, nil
}

func scanHabits(rows *sql.Rows) ([]domain.Habit, error) {
	var habits []domain.Habit
	for rows.Next() {
		var h domain.Habit
		var color, freq, createdAt string
		if err := rows.Scan(&h.ID, &h.Name, &h.PillarID, &color, &freq, &createdAt); err != nil {
			return nil, fmt.Errorf("escanear hábito: %w", err)
		}
		h.Color = domain.Color(color)
		h.Frequency = decodeFrequency(freq)
		h.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		habits = append(habits, h)
	}
	return habits, rows.Err()
}

func encodeFrequency(days []domain.Weekday) string {
	s := make([]string, len(days))
	for i, d := range days {
		s[i] = string(d)
	}
	return strings.Join(s, ",")
}

func decodeFrequency(s string) []domain.Weekday {
	parts := strings.Split(s, ",")
	days := make([]domain.Weekday, len(parts))
	for i, p := range parts {
		days[i] = domain.Weekday(p)
	}
	return days
}
