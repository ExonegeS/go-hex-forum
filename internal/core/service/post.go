package service

import (
	"context"
	"go-hex-forum/internal/core/domain"
)

type PostRepository interface {
	SavePost(ctx context.Context, post *domain.Post, userID int64) (int64, error)
	// GetPostByID(id int64) (*domain.Post, error)
	// GetRecentPosts(limit int) ([]domain.Post, error)
	// ArchiveOldPosts() error
}

type ImageStorage interface {
	UploadImage(bucket string, data []byte) (path string, err error)
	GetImageURL(path string) string
	// DeleteImage(path string) error
}

type PostService struct {
	PostRepo PostRepository
	// imageStorage ImageStorage
}

func NewPostService(PostRepo PostRepository) *PostService {
	return &PostService{PostRepo}

}

func (s *PostService) CreateNewPost(ctx context.Context, post *domain.Post) (int64, error) {
	userID := ctx.Value("user_id")
	return s.PostRepo.SavePost(ctx, post, userID.(int64))
}
