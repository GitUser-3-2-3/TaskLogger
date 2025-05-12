package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Category CategoryModel
	Tasks    TaskModel
	Session  SessionModel
	DB       *sql.DB
}

func NewModels(db *sql.DB) Models {
	return Models{
		Category: CategoryModel{DB: db},
		Tasks:    TaskModel{DB: db},
		Session:  SessionModel{DB: db},
		DB:       db,
	}
}
