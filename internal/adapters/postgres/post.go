package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-hex-forum/internal/core/domain"
)

type PostRepository struct {
	br BaseRepository
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{BaseRepository{db}}
}

func (r *PostRepository) SavePost(ctx context.Context, post *domain.Post, userID int64) (int64, error) {
	const op = "PostRepository.SavePost"

	query := `INSERT INTO posts (user_id, title, content, image_path, created_at, expires_at)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          RETURNING id`

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
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *PostRepository) GetActivePosts(ctx context.Context, pagination *domain.Pagination) ([]domain.Post, error) {
	const op = "PostRepository.GetActivePosts"

	var posts []domain.Post
	query := `SELECT p.id, u.id, u.name, COALESCE(u.avatar_url, '') AS avatar_url,
	                 p.title, p.content, COALESCE(p.image_path, '') AS image_path,
	                 p.created_at, p.is_archived 
	          FROM posts p
	          JOIN users u ON u.id = p.user_id
	          WHERE p.is_archived = false
	          ORDER BY p.created_at DESC
	          LIMIT $1 OFFSET $2`

	offset := (pagination.Page - 1) * pagination.PageSize

	rows, err := r.br.queryContext(ctx, query, pagination.PageSize, offset)
	if err != nil {
		return posts, fmt.Errorf("%s: queryContext: %w", op, err)
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
			return posts, fmt.Errorf("%s: rows.Scan: %w", op, err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return posts, fmt.Errorf("%s: rows.Err: %w", op, err)
	}

	return posts, nil
}

func (r *PostRepository) GetArchivedPosts(ctx context.Context, pagination *domain.Pagination) ([]domain.Post, error) {
	const op = "PostRepository.GetArchivedPosts"

	var posts []domain.Post
	query := `SELECT p.id, u.id, u.name, COALESCE(u.avatar_url, '') AS avatar_url,
	                 p.title, p.content, COALESCE(p.image_path, '') AS image_path,
	                 p.created_at, p.is_archived 
	          FROM posts p
	          JOIN users u ON u.id = p.user_id
	          WHERE p.is_archived = true
	          ORDER BY p.created_at DESC
	          LIMIT $1 OFFSET $2`

	offset := (pagination.Page - 1) * pagination.PageSize

	rows, err := r.br.queryContext(ctx, query, pagination.PageSize, offset)
	if err != nil {
		return posts, fmt.Errorf("%s: queryContext: %w", op, err)
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
			return posts, fmt.Errorf("%s: rows.Scan: %w", op, err)
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return posts, fmt.Errorf("%s: rows.Err: %w", op, err)
	}

	return posts, nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, postID int64) (domain.Post, error) {
	const op = "PostRepository.GetPostByID"

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
		return post, fmt.Errorf("%s: %w", op, err)
	}

	return post, nil
}

func (r *PostRepository) UpdateExpiresAt(ctx context.Context, postID int64, date time.Time) error {
	const op = "PostRepository.UpdateExpiresAt"

	query := `UPDATE posts SET expires_at = $1 WHERE id = $2`

	_, err := r.br.execContext(ctx, query, date, postID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PostRepository) ArchiveExpiredPosts(ctx context.Context) error {
	const op = "PostRepository.ArchiveExpiredPosts"

	query := `UPDATE posts SET is_archived = true WHERE expires_at <= NOW() AND is_archived = false`

	_, err := r.br.execContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
