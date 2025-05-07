package data

import "database/sql"

type Models struct {
	Tasks    TaskModel
	Category CategoryModel
}

func NewModels(db *sql.DB) *Models {
	return &Models{
		Tasks:    TaskModel{DB: db},
		Category: CategoryModel{DB: db},
	}
}
