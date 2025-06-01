package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Session struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
	Duration  int       `json:"duration"`
	Note      string    `json:"note"`
}

type SessionModel struct {
	DB *sql.DB
}

func (dbm *SessionModel) InsertTx(tx *sql.Tx, session *Session) error {
	session.ID = uuid.New().String()

	query := `INSERT INTO sessions (session_id, task_id, started_at, ended_at, duration, notes)
		   VALUES ($1, $2, $3, $4, $5, $6)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		session.ID, session.TaskID, session.StartedAt, session.EndedAt,
		session.Duration, session.Note,
	}
	row := tx.QueryRowContext(ctx, query, args...)
	if err := row.Err(); err != nil {
		return row.Err()
	}
	return nil
}

func (dbm *SessionModel) GetForTask(taskId int64) ([]*Session, error) {
	query := `SELECT * FROM sessions WHERE task_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := dbm.DB.QueryContext(ctx, query, taskId)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("couldn't close rows")
		}
	}()
	var sessions []*Session

	for rows.Next() {
		var session Session
		err = rows.Scan(&session.ID,
			&session.TaskID,
			&session.StartedAt,
			&session.EndedAt,
			&session.Duration,
			&session.Note,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return sessions, nil
}
