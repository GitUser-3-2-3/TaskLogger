package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PriorityType string
type StatusType string

const (
	StatusPending    StatusType = "Pending"
	StatusInProgress StatusType = "In Progress"
	StatusPaused     StatusType = "Paused"
	StatusCompleted  StatusType = "Completed"
	StatusCancelled  StatusType = "Cancelled"
)

const (
	PriorityLow    PriorityType = "low"
	PriorityMedium PriorityType = "medium"
	PriorityHigh   PriorityType = "high"
	PriorityUrgent PriorityType = "urgent"
)

type Tasks struct {
	ID          string       `json:"task_id"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Status      StatusType   `json:"status"`
	Priority    PriorityType `json:"priority"`
	ImageUrl    string       `json:"image_url,omitempty"`
	Duration    int          `json:"duration"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Deadline    *time.Time   `json:"deadline,omitempty"`
	UserID      string       `json:"-"`
	CategoryID  *int64       `json:"category_id,omitempty"`
}

type TaskModel struct {
	DB *sql.DB
}

func (dbm *TaskModel) Insert(task *Tasks) error {
	task.ID = uuid.New().String()

	query := `INSERT INTO tasks (task_id, name, description, status, priority, image_url, deadline, 
		   user_id, category_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{
		task.ID, task.Name, task.Description, task.Status, task.Priority, task.ImageUrl,
		task.Deadline, task.UserID, task.CategoryID,
	}
	row := dbm.DB.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}

func (dbm *TaskModel) GetByTaskId(taskId string) (*Tasks, error) {
	if taskId == "" || !isValidUUID(taskId) {
		return nil, ErrRecordNotFound
	}
	query := `SELECT task_id, name, description, status, priority, image_url, duration_minutes, 
		   created_at, updated_at, deadline, user_id, category_id FROM tasks WHERE task_id = $1`
	var task Tasks

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := dbm.DB.QueryRowContext(ctx, query, taskId).Scan(&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.ImageUrl,
		&task.Duration,
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
	query := `SELECT task_id, name, description, status, priority, image_url, duration_minutes,
                created_at, updated_at, deadline, user_id, category_id 
                FROM tasks WHERE category_id = $1`

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
			&task.ImageUrl,
			&task.Duration,
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
	if task.ID == "" || !isValidUUID(task.ID) {
		return ErrRecordNotFound
	}
	query := `UPDATE tasks SET name = $1, description = $2, status = $3, image_url = $4,
                 duration_minutes = $5, priority = $6, deadline = $7, category_id = $8
                 WHERE task_id = $9`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{
		task.Name, task.Description, task.Status, task.ImageUrl, task.Duration, task.Priority,
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

func (dbm *TaskModel) Delete(taskId string) error {
	if taskId == "" || !isValidUUID(taskId) {
		return ErrRecordNotFound
	}
	query := `DELETE FROM tasks WHERE task_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := dbm.DB.ExecContext(ctx, query, taskId)
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

// Transactional Methods

func (dbm *TaskModel) GetByTaskIdTx(tx *sql.Tx, taskId string) (*Tasks, error) {
	if taskId == "" || !isValidUUID(taskId) {
		return nil, ErrRecordNotFound
	}
	query := `SELECT task_id, name, description, status, priority, image_url, duration_minutes, 
		   created_at, updated_at, deadline, user_id, category_id 
		   FROM tasks WHERE task_id = $1`

	var task Tasks

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, taskId).Scan(&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.ImageUrl,
		&task.Duration,
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

func (dbm *TaskModel) UpdateTx(tx *sql.Tx, task *Tasks) error {
	if task.ID == "" || !isValidUUID(task.ID) {
		return ErrRecordNotFound
	}
	query := `UPDATE tasks SET name = $1, description = $2, status = $3, image_url = $4, 
                 duration_minutes = $5, priority = $6, deadline = $7, category_id = $8 
                 WHERE task_id = $9`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		task.Name, task.Description, task.Status, task.ImageUrl, task.Duration, task.Priority,
		task.Deadline, task.CategoryID, task.ID,
	}
	result, err := tx.ExecContext(ctx, query, args...)
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
