package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("abrir banco: %w", err)
	}
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrar banco: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS pillars (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT    NOT NULL UNIQUE,
			created_at DATETIME NOT NULL
		);

		CREATE TABLE IF NOT EXISTS habits (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT    NOT NULL,
			pillar_id  INTEGER NOT NULL REFERENCES pillars(id),
			color      TEXT    NOT NULL,
			frequency  TEXT    NOT NULL,
			created_at DATETIME NOT NULL,
			UNIQUE(name, pillar_id)
		);

		CREATE TABLE IF NOT EXISTS tasks (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			title       TEXT    NOT NULL UNIQUE,
			description TEXT    NOT NULL DEFAULT '',
			status      TEXT    NOT NULL DEFAULT 'em_andamento',
			next_action TEXT,
			created_at  DATETIME NOT NULL,
			updated_at  DATETIME
		);

		CREATE TABLE IF NOT EXISTS sessions (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			ref_type    TEXT    NOT NULL,
			ref_id      INTEGER NOT NULL,
			started_at  DATETIME NOT NULL,
			finished_at DATETIME
		);
	`)
	return err
}
