package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"habits/internal/domain"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Save(t domain.Task) (domain.Task, error) {
	res, err := r.db.Exec(
		`INSERT INTO tasks (title, description, status, next_action, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		t.Title, t.Description, string(t.Status), encodeNextAction(t.NextAction),
		t.CreatedAt.UTC().Format(time.RFC3339), encodeTime(t.UpdatedAt),
	)
	if err != nil {
		return domain.Task{}, fmt.Errorf("inserir tarefa: %w", err)
	}
	id, _ := res.LastInsertId()
	t.ID = id
	return t, nil
}

func (r *TaskRepository) FindByTitle(title string) (domain.Task, bool, error) {
	row := r.db.QueryRow(`SELECT id, title, description, status, next_action, created_at, updated_at FROM tasks WHERE title = ?`, title)
	return scanTask(row)
}

func (r *TaskRepository) FindInProgress() ([]domain.Task, error) {
	rows, err := r.db.Query(`SELECT id, title, description, status, next_action, created_at, updated_at FROM tasks WHERE status = 'em_andamento'`)
	if err != nil {
		return nil, fmt.Errorf("listar tarefas: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var t domain.Task
		var status, nextAction, createdAt string
		var updatedAt sql.NullString
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &status, &nextAction, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("escanear tarefa: %w", err)
		}
		t.Status = domain.TaskStatus(status)
		t.NextAction = decodeNextAction(nextAction)
		t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		if updatedAt.Valid && updatedAt.String != "" {
			tt, _ := time.Parse(time.RFC3339, updatedAt.String)
			t.UpdatedAt = &tt
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) Update(t domain.Task) error {
	_, err := r.db.Exec(
		`UPDATE tasks SET title=?, description=?, status=?, next_action=?, updated_at=? WHERE id=?`,
		t.Title, t.Description, string(t.Status), encodeNextAction(t.NextAction),
		encodeTime(t.UpdatedAt), t.ID,
	)
	return err
}

func (r *TaskRepository) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	return err
}

func scanTask(row *sql.Row) (domain.Task, bool, error) {
	var t domain.Task
	var status, nextAction, createdAt string
	var updatedAt sql.NullString
	err := row.Scan(&t.ID, &t.Title, &t.Description, &status, &nextAction, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return domain.Task{}, false, nil
	}
	if err != nil {
		return domain.Task{}, false, fmt.Errorf("escanear tarefa: %w", err)
	}
	t.Status = domain.TaskStatus(status)
	t.NextAction = decodeNextAction(nextAction)
	t.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	if updatedAt.Valid && updatedAt.String != "" {
		tt, _ := time.Parse(time.RFC3339, updatedAt.String)
		t.UpdatedAt = &tt
	}
	return t, true, nil
}

func encodeNextAction(na *domain.NextAction) string {
	if na == nil {
		return ""
	}
	return string(na.Type)
}

func decodeNextAction(s string) *domain.NextAction {
	if s == "" {
		return nil
	}
	return &domain.NextAction{Type: domain.NextActionType(s)}
}

func encodeTime(t *time.Time) sql.NullString {
	if t == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: t.UTC().Format(time.RFC3339), Valid: true}
}
