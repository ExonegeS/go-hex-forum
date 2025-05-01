package service

import (
	"context"
	"go-hex-forum/internal/core/domain"
	"log/slog"
	"time"
)

type CommentRepository interface {
	// Основные CRUD операции
	Save(ctx context.Context, comment *domain.Comment) error
	GetByID(ctx context.Context, id int64) (*domain.Comment, error)
	GetByPostID(ctx context.Context, postID string, limit int) ([]domain.Comment, error)

	// Для управления жизненным циклом
	DeleteCommentsForPost(ctx context.Context, postID string) error
	UpdateCommentContent(ctx context.Context, id, newContent string) error

	// Для связей между комментариями
	GetReplies(ctx context.Context, parentID string) ([]domain.Comment, error)
	GetThread(ctx context.Context, rootID string) ([]domain.Comment, error)

	// Для обновления времени поста
	GetLastCommentTimeForPost(ctx context.Context, postID string) (time.Time, error)
}

type CommentService struct {
	commentRepo  CommentRepository
	postRepo     PostRepository
	imageStorage ImageStorage
}

func NewCommentService(cr CommentRepository, pr PostRepository,
	is ImageStorage, logger *slog.Logger) *CommentService {
	return &CommentService{
		commentRepo:  cr,
		postRepo:     pr,
		imageStorage: is,
	}
}
