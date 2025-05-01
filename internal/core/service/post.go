package service

import (
	"go-hex-forum/internal/core/domain"
)

type PostRepository interface {
	SavePost(post *domain.Post) error
	GetPostByID(id int64) (*domain.Post, error)
	GetRecentPosts(limit int) ([]domain.Post, error)
	ArchiveOldPosts() error
}

type ImageStorage interface {
	UploadImage(bucket string, data []byte) (path string, err error)
	GetImageURL(path string) string
	// DeleteImage(path string) error
}

type PostService struct {
	repo         PostRepository
	imageStorage ImageStorage
}
