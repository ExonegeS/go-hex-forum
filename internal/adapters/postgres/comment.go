package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go-hex-forum/internal/core/domain"
)

type CommentRepository struct {
	br BaseRepository
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{BaseRepository{db}}
}

func (r *CommentRepository) SaveComment(ctx context.Context, comment *domain.Comment) (int64, error) {
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
	err := r.br.queryRowContext(ctx, query,
		comment.PostID,
		comment.Author.ID,
		comment.Content,
		imagePath,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to save comment: %w", err)
	}

	comment.ID = id
	return id, nil
}

func (r *CommentRepository) GetByPostID(ctx context.Context, postID int64) ([]*domain.Comment, error) {
	query := `
        SELECT c.id, u.name, COALESCE(u.avatar_url, '') AS avatar_url, 
               c.content, c.image_path, c.created_at
        FROM comments c
        JOIN users u ON u.id = c.user_id
        WHERE c.post_id = $1
        ORDER BY c.created_at ASC
    `

	rows, err := r.br.queryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("querying comments: %w", err)
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		var c domain.Comment
		var imagePath sql.NullString
		err := rows.Scan(&c.ID, &c.Author.Name, &c.Author.AvatarURL, &c.Content, &imagePath, &c.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning comment: %w", err)
		}
		if imagePath.Valid {
			c.ImagePath = imagePath.String
		} else {
			c.ImagePath = ""
		}
		comments = append(comments, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return comments, nil
}
