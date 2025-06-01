package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	ErrDuplicateEntry = errors.New("category already exists")
)

type Categories struct {
	ID        int64     `json:"category_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	UserID    string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type CategoryModel struct {
	DB *sql.DB
}

func (dbm *CategoryModel) Insert(ctg *Categories) error {
	query := `INSERT INTO categories (name, color, user_id) VALUES ($1, $2, $3) 
		   RETURNING ctg_id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{ctg.Name, ctg.Color, ctg.UserID}
	err := dbm.DB.QueryRowContext(ctx, query, args...).Scan(&ctg.ID)

	var pgErr *pq.Error
	if err != nil {
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return ErrDuplicateEntry
		default:
			return err
		}
	}
	return nil
}

func (dbm *CategoryModel) GetByCtgId(ctgId int64) (*Categories, error) {
	if ctgId < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT ctg_id, name, color, user_id, created_at FROM categories 
                 WHERE ctg_id = $1`
	var ctg Categories

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := dbm.DB.QueryRowContext(ctx, query, ctgId).Scan(&ctg.ID,
		&ctg.Name,
		&ctg.Color,
		&ctg.UserID,
		&ctg.CreatedAt,
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
	if ctg.ID < 1 {
		return ErrRecordNotFound
	}
	query := `UPDATE categories SET name = $2, color = $3 WHERE ctg_id = $1
                 RETURNING ctg_id, name, color`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{ctg.ID, ctg.Name, ctg.Color}
	err := dbm.DB.QueryRowContext(ctx, query, args...).Scan(&ctg.ID,
		&ctg.Name,
		&ctg.Color,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (dbm *CategoryModel) Delete(ctgId int64) error {
	if ctgId < 1 {
		return ErrRecordNotFound
	}
	query := `DELETE FROM categories WHERE ctg_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := dbm.DB.ExecContext(ctx, query, ctgId)
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

func (dbm *CategoryModel) GetAllByUserId(userId string) ([]*Categories, error) {
	if userId == "" || !isValidUUID(userId) {
		return nil, ErrRecordNotFound
	}
	query := `SELECT ctg_id, name, color, created_at FROM categories WHERE user_id = $1`

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
		err = rows.Scan(&ctg.ID, &ctg.Name, &ctg.Color, &ctg.CreatedAt)
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
