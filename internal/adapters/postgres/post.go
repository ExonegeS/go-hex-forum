package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"time"
)

type PostRepository struct {
	br BaseRepository
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{BaseRepository{db}}
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
	err := r.br.queryRowContext(ctx, query,
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

func (r *PostRepository) GetActivePosts(ctx context.Context, pagination *domain.Pagination) ([]domain.Post, error) {
	var posts []domain.Post
	query := `SELECT p.id,u.id,u.name,COALESCE(u.avatar_url, '') AS avatar_url,p.title,p.content,COALESCE(p.image_path, '') AS image_path,p.created_at,p.is_archived 
	FROM posts p
	JOIN users u ON u.id = p.user_id
	WHERE p.is_archived = false
	ORDER BY p.created_at DESC
	LIMIT $1 OFFSET $2`

	offset := (pagination.Page - 1) * pagination.PageSize

	rows, err := r.br.queryContext(ctx, query, pagination.PageSize, offset)
	if err != nil {
		return posts, fmt.Errorf("failed to query active posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post domain.Post
		err := rows.Scan(
			&post.ID,
			&post.PostAuthor.ID,
			&post.PostAuthor.Name,
			&post.PostAuthor.AvatarURL,
			&post.Title,
			&post.Content,
			&post.ImagePath,
			&post.CreatedAt,
			&post.IsArchived,
		)
		if err != nil {
			return posts, fmt.Errorf("failed to query active posts: %w", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return posts, fmt.Errorf("failed to query active posts: %w", err)
	}
	return posts, nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, postID int64) (domain.Post, error) {
	var post domain.Post

	query := `
    SELECT 
        p.id,
        u.id,
        u.name,
        COALESCE(u.avatar_url, '') AS avatar_url,
        p.title,
        p.content,
        COALESCE(p.image_path, '') AS image_path,
        p.created_at,
		p.expires_at,
        p.is_archived
    FROM posts p
    JOIN users u ON u.id = p.user_id
    WHERE p.id = $1
    `
	err := r.br.queryRowContext(ctx, query, postID).Scan(
		&post.ID,
		&post.PostAuthor.ID,
		&post.PostAuthor.Name,
		&post.PostAuthor.AvatarURL,
		&post.Title,
		&post.Content,
		&post.ImagePath,
		&post.CreatedAt,
		&post.ExpiresAt,
		&post.IsArchived,
	)
	if err != nil {
		return post, fmt.Errorf("GetPostByID: %w", err)
	}
	return post, nil
}

func (r *PostRepository) UpdateExpiresAt(ctx context.Context, postID int64, date time.Time) error {
	query := `
		UPDATE posts
		SET expires_at = $1
		WHERE id = $2
	`

	_, err := r.br.execContext(ctx, query, date, postID)
	if err != nil {
		return fmt.Errorf("UpdateExpiresAt: %w", err)
	}
	return nil
}
