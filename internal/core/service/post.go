package service

import (
	"context"
	"errors"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"time"
)

type PostRepository interface {
	SavePost(ctx context.Context, post *domain.Post, userID int64) (int64, error)
	GetActivePosts(ctx context.Context, pagination *domain.Pagination) ([]domain.Post, error)
	GetPostByID(ctx context.Context, postID int64) (domain.Post, error)
}

type ImageStorage interface {
	UploadImage(ctx context.Context, userID int64, data []byte) (publicURL string, err error)
	GetImageURL(userID int64, code string) string
}

type PostService struct {
	postRepo     PostRepository
	imageStorage ImageStorage
}

func NewPostService(postRepo PostRepository, imageStorage ImageStorage) *PostService {
	return &PostService{
		postRepo,
		imageStorage,
	}
}

func (s *PostService) CreateNewPost(ctx context.Context, post *domain.Post, imageData []byte) (int64, error) {
	if post.Title == "" || post.Content == "" {
		return 0, errors.New("title and content are required")
	}
	post.CreatedAt = time.Now().UTC()

	if len(imageData) > 0 {
		url, err := s.imageStorage.UploadImage(ctx, post.PostAuthor.ID, imageData)
		if err != nil {
			return -1, err
		}
		post.ImagePath = url
	}

	return s.postRepo.SavePost(ctx, post, post.PostAuthor.ID)
}

func (s *PostService) GetActivePosts(ctx context.Context) ([]domain.Post, error) {
	pagination := &domain.Pagination{
		Page:     1,
		PageSize: 10,
	}
	return s.postRepo.GetActivePosts(ctx, pagination)
}

func (s *PostService) GetPostByID(ctx context.Context, postID int64) (domain.Post, error) {
	return s.postRepo.GetPostByID(ctx, postID)
}

func (s *PostService) UploadImage(ctx context.Context, userID int64, imageData []byte) (string, error) {
	publicURL, err := s.imageStorage.UploadImage(ctx, userID, imageData)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}
	return publicURL, nil
}
