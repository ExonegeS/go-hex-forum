package service

import (
	"context"
	"go-hex-forum/internal/core/domain"
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

type CommentService struct {
	commentRepo  CommentRepository
	postRepo     PostRepository
	imageStorage ImageStorage
}

func NewCommentService(cr CommentRepository, pr PostRepository,
	is ImageStorage) *CommentService {
	return &CommentService{
		commentRepo:  cr,
		postRepo:     pr,
		imageStorage: is,
	}
}

func (s *CommentService) SaveComment(ctx context.Context, comment *domain.Comment, imageData []byte) (int64, error) {
	if len(imageData) > 0 {
		url, err := s.imageStorage.UploadImage(ctx, comment.Author.ID, imageData)
		if err != nil {
			return -1, err
		}
		comment.ImagePath = url
	}
	return s.commentRepo.SaveComment(ctx, comment)
}

func (s *CommentService) GetByPostID(ctx context.Context, postID int64) ([]*domain.Comment, error) {
	return s.commentRepo.GetByPostID(ctx, postID)
}
