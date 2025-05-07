package data

import (
	"database/sql"
	"time"
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
	CategoryID    *int64     `json:"-"`
}

type TaskModel struct {
	DB *sql.DB
}

func (ctg *TaskModel) Insert(task *Tasks) error {
	return nil
}

func (ctg *TaskModel) Get(id int64) (*Tasks, error) {
	return nil, nil
}

func (ctg *TaskModel) Update(task *Tasks) error {
	return nil
}

func (ctg *TaskModel) Delete(id int64) error {
	return nil
}
