package data

import "database/sql"

type Categories struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	UserID int64  `json:"-"`
}

type CategoriesModel struct {
	DB *sql.DB
}
