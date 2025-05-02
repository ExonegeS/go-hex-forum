package service

import (
	"context"
	"errors"
	"fmt"
	"go-hex-forum/internal/core/domain"
)

type PostRepository interface {
	SavePost(ctx context.Context, post *domain.Post, userID int64) (int64, error)
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

func (s *PostService) CreateNewPost(
	ctx context.Context,
	title, content string,
	imageData []byte,
	userID int64,
) (int64, error) {
	if title == "" || content == "" {
		return 0, errors.New("title and content are required")
	}

	var imagePath string
	// if len(imageData) > 0 {
	// 	path, err := s.imageStorage.UploadImage(ctx, userID, imageData)
	// 	if err != nil {
	// 		return 0, fmt.Errorf("failed to upload image: %w", err)
	// 	}
	// 	imagePath = s.imageStorage.GetImageURL(path)
	// }

	post := &domain.Post{
		Title:     title,
		Content:   content,
		ImagePath: imagePath,
	}

	return s.postRepo.SavePost(ctx, post, userID)
}

func (s *PostService) UploadImage(ctx context.Context, userID int64, imageData []byte) (string, error) {
	publicURL, err := s.imageStorage.UploadImage(ctx, userID, imageData)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}
	return publicURL, nil
}
