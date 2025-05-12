package data

import (
	"context"
	"database/sql"
	"time"
)

type SessionType string

const (
	SessionWork   SessionType = "work"
	SessionsBreak SessionType = "break"
)

type Session struct {
	ID           int64       `json:"id"`
	TaskID       int64       `json:"task_id"`
	SessionStart time.Time   `json:"session_start"`
	SessionEnd   time.Time   `json:"session_end"`
	Duration     int         `json:"duration"`
	Note         string      `json:"note"`
	SessionType  SessionType `json:"session_type"`
}

type SessionModel struct {
	DB *sql.DB
}

func (dbm *SessionModel) InsertTx(tx *sql.Tx, session *Session) (int64, error) {
	query := `INSERT INTO sessions (task_id, session_start, session_end, duration, note, session_type)
		    VALUES (?, ?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		session.TaskID, session.SessionStart, session.SessionEnd, session.Duration,
		session.Note, session.SessionType,
	}
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
