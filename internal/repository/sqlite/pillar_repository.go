package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"habits/internal/domain"
)

type PillarRepository struct {
	db *sql.DB
}

func NewPillarRepository(db *sql.DB) *PillarRepository {
	return &PillarRepository{db: db}
}

func (r *PillarRepository) Save(p domain.Pillar) (domain.Pillar, error) {
	res, err := r.db.Exec(
		`INSERT INTO pillars (name, created_at) VALUES (?, ?)`,
		p.Name, p.CreatedAt.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return domain.Pillar{}, fmt.Errorf("inserir pilar: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Pillar{}, fmt.Errorf("obter id: %w", err)
	}
	p.ID = id
	return p, nil
}

func (r *PillarRepository) FindByID(id int64) (domain.Pillar, bool, error) {
	row := r.db.QueryRow(`SELECT id, name, created_at FROM pillars WHERE id = ?`, id)
	return scanPillar(row)
}

func (r *PillarRepository) FindByName(name string) (domain.Pillar, bool, error) {
	row := r.db.QueryRow(`SELECT id, name, created_at FROM pillars WHERE name = ?`, name)
	return scanPillar(row)
}

func (r *PillarRepository) FindAll() ([]domain.Pillar, error) {
	rows, err := r.db.Query(`SELECT id, name, created_at FROM pillars`)
	if err != nil {
		return nil, fmt.Errorf("listar pilares: %w", err)
	}
	defer rows.Close()

	var pillars []domain.Pillar
	for rows.Next() {
		var p domain.Pillar
		var createdAt string
		if err := rows.Scan(&p.ID, &p.Name, &createdAt); err != nil {
			return nil, fmt.Errorf("escanear pilar: %w", err)
		}
		p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		pillars = append(pillars, p)
	}
	return pillars, rows.Err()
}

func (r *PillarRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM pillars WHERE id = ?`, id)
	return err
}

func scanPillar(row *sql.Row) (domain.Pillar, bool, error) {
	var p domain.Pillar
	var createdAt string
	err := row.Scan(&p.ID, &p.Name, &createdAt)
	if err == sql.ErrNoRows {
		return domain.Pillar{}, false, nil
	}
	if err != nil {
		return domain.Pillar{}, false, fmt.Errorf("escanear pilar: %w", err)
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	return p, true, nil
}
