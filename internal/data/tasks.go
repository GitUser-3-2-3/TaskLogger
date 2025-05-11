package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

type StatusType string

const (
	StatusNotStarted StatusType = "Not Started"
	StatusInProgress StatusType = "In Progress"
	StatusPaused     StatusType = "Paused"
	StatusCompleted  StatusType = "Completed"
)

type Tasks struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description,omitempty"`
	Status        StatusType `json:"status"`
	Priority      int        `json:"priority"`
	Image         string     `json:"image,omitempty"`
	TotalDuration int        `json:"total-duration"`
	CreatedAt     time.Time  `json:"created-at"`
	UpdatedAt     time.Time  `json:"updated-at"`
	Deadline      *time.Time `json:"deadline,omitempty"`
	UserID        int64      `json:"-"`
	CategoryID    *int64     `json:"category_id,omitempty"`
}

type TaskModel struct {
	DB *sql.DB
}

func (dbm *TaskModel) Insert(task *Tasks) (int64, error) {
	query := `INSERT INTO tasks (name, description, status, priority, image, deadline, user_id, category_id) 
		    VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		task.Name, task.Description, task.Status, task.Priority, task.Image,
		task.Deadline, task.UserID, task.CategoryID,
	}
	rows, err := dbm.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	id, err := rows.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (dbm *TaskModel) GetById(id int64) (*Tasks, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, name, description, status, priority, image, total_duration, 
		    created_at, updated_at, deadline, tasks.user_id, category_id FROM tasks WHERE id = ?`
	var task Tasks

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := dbm.DB.QueryRowContext(ctx, query, id).Scan(&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.Image,
		&task.TotalDuration,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.Deadline,
		&task.UserID,
		&task.CategoryID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &task, nil
}

func (dbm *TaskModel) GetAllByCategory(ctgID int64) ([]*Tasks, error) {
	if ctgID < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, name, description, status, priority, image, total_duration,
                created_at, updated_at, deadline, user_id, category_id 
                FROM tasks WHERE category_id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbm.DB.QueryContext(ctx, query, ctgID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close rows")
		}
	}()
	var tasks []*Tasks

	for rows.Next() {
		var task Tasks
		err = rows.Scan(&task.ID,
			&task.Name,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.Image,
			&task.TotalDuration,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Deadline,
			&task.UserID,
			&task.CategoryID,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (dbm *TaskModel) Update(task *Tasks) error {
	query := `UPDATE tasks SET name = ?, description = ?, status = ?, image = ?, 
		    priority = ?, deadline = ?, category_id = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		task.Name, task.Description, task.Status, task.Image, task.Priority,
		task.Deadline, task.CategoryID, task.ID,
	}
	result, err := dbm.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (dbm *TaskModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM tasks WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := dbm.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrRecordNotFound
	}
	return nil
}
