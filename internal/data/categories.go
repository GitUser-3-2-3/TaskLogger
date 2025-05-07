package data

import (
	"TaskLogger/internal/validator"
	"database/sql"
)

type Categories struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	UserID int64  `json:"-"`
}

type CategoriesModel struct {
	DB *sql.DB
}

func ValidateCategory(vld *validator.Validator, category *Categories) {
	vld.CheckError(category.Name != "", "name", "must not be empty")
	vld.CheckError(len(category.Name) > 0 &&
		len(category.Name) <= 50, "name", "cannot be longer than 50 chars")

	vld.CheckError(category.UserID > 0, "user_id", "cannot be zero or negative")
}
