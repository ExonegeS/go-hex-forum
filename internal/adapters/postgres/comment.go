package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go-hex-forum/internal/core/domain"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db}
}

func (r *CommentRepository) SaveNewComment(ctx context.Context, comment *domain.Comment, userID int64) (int64, error) {
	query := `
        INSERT INTO comments (post_id, user_id, content, image_path)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	// Handle nullable image path
	imagePath := sql.NullString{
		String: comment.ImagePath,
		Valid:  comment.ImagePath != "",
	}

	var id int64
	err := r.db.QueryRowContext(ctx, query,
		comment.PostID,
		userID,
		comment.Content,
		imagePath,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to save comment: %w", err)
	}

	comment.ID = id
	return id, nil
}
