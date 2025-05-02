package postgres

import (
	"context"
	"database/sql"
	"go-hex-forum/internal/core/domain"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db}
}

func (r *PostRepository) SavePost(ctx context.Context, post *domain.Post, userID int64) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx, `
	INSERT INTO posts (user_id, title, content, created_at, expires_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`,
		userID,
		post.Title,
		post.Content,
		post.CreatedAt,
		post.ExpiresAt,
	).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}
