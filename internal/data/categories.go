package data

import (
	"TaskLogger/internal/validator"
	"context"
	"database/sql"
	"time"
)

type Categories struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	UserID int64  `json:"-"`
}

type CategoryModel struct {
	DB *sql.DB
}

func (dbm *CategoryModel) Insert(ctg *Categories) (int64, error) {
	query := `INSERT INTO categories (id, name, color, user_id)
                VALUES (?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{ctg.ID, ctg.Name, ctg.Color, ctg.UserID}
	result, err := dbm.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (dbm *CategoryModel) Get(id int64) (*Categories, error) {
	return nil, nil
}

func (dbm *CategoryModel) Update(category *Categories) error {
	return nil
}

func (dbm *CategoryModel) Delete(id int64) error {
	return nil
}

func ValidateCategory(vld *validator.Validator, category *Categories) {
	vld.CheckError(category.Name != "", "name", "must not be empty")
	vld.CheckError(len(category.Name) > 0 &&
		len(category.Name) <= 50, "name", "cannot be longer than 50 chars")

	vld.CheckError(category.UserID > 0, "user_id", "cannot be zero or negative")
}
