package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ExonegeS/go-hex-forum/internal/adapter/outbound/repository/postgres/entity"
	"github.com/ExonegeS/go-hex-forum/internal/domain/model"
)

type XRepository struct {
	db *sql.DB
}

func NewXRepository(db *sql.DB) *XRepository {
	return &XRepository{db: db}
}

func (r *XRepository) Save(ctx context.Context, x model.X) (int, error) {
	op := "XRepository.Save"
	e := entity.XToEntity(&x)

	query := "INSERT INTO x (data) VALUES ($1) RETURNING id"
	row := r.db.QueryRowContext(ctx, query, e.Data)

	var id int
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil
		}
		err = fmt.Errorf("%s: %w", op, err)
		return -1, err
	}

	return id, nil
}
