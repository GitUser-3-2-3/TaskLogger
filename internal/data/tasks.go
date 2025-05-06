package data

import "time"

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
	Priority      byte       `json:"priority"`
	Image         string     `json:"image,omitempty"`
	TotalDuration int        `json:"total-duration"`
	CreatedAt     time.Time  `json:"created-at"`
	UpdatedAt     time.Time  `json:"updated-at"`
	Deadline      time.Time  `json:"deadline"`
	UserID        int64      `json:"-"`
	CategoryID    int64      `json:"-"`
}
