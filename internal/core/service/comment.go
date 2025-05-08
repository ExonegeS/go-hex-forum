package service

import (
	"context"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"time"
)

type CommentRepository interface {
	// Основные CRUD операции
	SaveComment(ctx context.Context, comment *domain.Comment) (int64, error)
	// GetByID(ctx context.Context, id int64) (*domain.Comment, error)
	GetByPostID(ctx context.Context, postID int64) ([]*domain.Comment, error)

	// // Для управления жизненным циклом
	// DeleteCommentsForPost(ctx context.Context, postID string) error
	// UpdateCommentContent(ctx context.Context, id, newContent string) error

	// // Для связей между комментариями
	// GetReplies(ctx context.Context, parentID string) ([]domain.Comment, error)
	// GetThread(ctx context.Context, rootID string) ([]domain.Comment, error)

	// // Для обновления времени поста
	// GetLastCommentTimeForPost(ctx context.Context, postID string) (time.Time, error)
}

// type IPostExpire interface // get post | update post expire at

type CommentPostRepo interface {
	GetPostByID(ctx context.Context, postID int64) (domain.Post, error)
	UpdateExpiresAt(ctx context.Context, postID int64, timeInSec time.Time) error
}

type CommentService struct {
	transactor   Transactor
	commentRepo  CommentRepository
	postRepo     CommentPostRepo
	imageStorage ImageStorage
}

func NewCommentService(tr Transactor, cr CommentRepository, pr CommentPostRepo,
	is ImageStorage) *CommentService {
	return &CommentService{
		transactor:   tr,
		commentRepo:  cr,
		postRepo:     pr,
		imageStorage: is,
	}
}

func (s *CommentService) SaveComment(ctx context.Context, comment *domain.Comment, imageData []byte) (int64, error) {
	post, err := s.postRepo.GetPostByID(ctx, comment.PostID)
	if err != nil {
		return -1, fmt.Errorf("post not found: %w", err)
	}
	if post.IsArchived {
		return -1, fmt.Errorf("post is archived, new comments are prohibited")
	}
	if len(imageData) > 0 {
		url, err := s.imageStorage.UploadImage(ctx, comment.Author.ID, imageData)
		if err != nil {
			return -1, err
		}
		comment.ImagePath = url
	}
	var id int64 = -1
	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if id, err = s.commentRepo.SaveComment(txCtx, comment); err != nil {
			return fmt.Errorf("failed to save comment: %w", err)
		}

		if err := s.postRepo.UpdateExpiresAt(txCtx, comment.PostID, time.Now().Add(15*time.Minute)); err != nil {
			return fmt.Errorf("failed to update post: %w", err)
		}

		return nil
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *CommentService) GetByPostID(ctx context.Context, postID int64) ([]*domain.Comment, error) {
	return s.commentRepo.GetByPostID(ctx, postID)
}
