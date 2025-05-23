package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

var (
	ErrDuplicateEntry = errors.New("category already exists")
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
	query := `INSERT INTO categories (name, color, user_id)
                VALUES (?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{ctg.Name, ctg.Color, ctg.UserID}
	result, err := dbm.DB.ExecContext(ctx, query, args...)

	var mysqlErr *mysql.MySQLError
	if err != nil {
		switch {
		case errors.As(err, &mysqlErr) && mysqlErr.Number == 1062:
			return 0, ErrDuplicateEntry
		default:
			return 0, err
		}
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (dbm *CategoryModel) GetById(id int64) (*Categories, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, name, color, user_id FROM categories WHERE id = ?`
	var ctg Categories

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := dbm.DB.QueryRowContext(ctx, query, id).Scan(&ctg.ID,
		&ctg.Name,
		&ctg.Color,
		&ctg.UserID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &ctg, nil
}

func (dbm *CategoryModel) Update(ctg *Categories) error {
	query := `UPDATE categories SET name = ?, color = ? WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := dbm.DB.ExecContext(ctx, query, ctg.Name, ctg.Color, ctg.ID)
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

func (dbm *CategoryModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM categories WHERE id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := dbm.DB.ExecContext(ctx, query, id)
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

func (dbm *CategoryModel) GetAllByUserId(userId int64) ([]*Categories, error) {
	if userId < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, name, color FROM categories WHERE user_id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := dbm.DB.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close rows")
		}
	}(rows)
	var ctgList []*Categories

	for rows.Next() {
		var ctg Categories
		err = rows.Scan(&ctg.ID, &ctg.Name, &ctg.Color)
		if err != nil {
			return nil, err
		}
		ctgList = append(ctgList, &ctg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ctgList, nil
}
