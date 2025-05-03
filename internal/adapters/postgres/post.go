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
	query := `INSERT INTO posts (user_id, title, content,image_path, created_at, expires_at)
	VALUES ($1, $2, $3, $4, $5,$6)
	RETURNING id`
	// Handle nullable image path
	imagePath := sql.NullString{
		String: post.ImagePath,
		Valid:  post.ImagePath != "",
	}
	var id int64
	err := r.db.QueryRowContext(ctx, query,
		userID,
		post.Title,
		post.Content,
		imagePath,
		post.CreatedAt,
		post.ExpiresAt,
	).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}
