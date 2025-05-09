package service

import (
	"context"
	"fmt"
	"time"

	"go-hex-forum/internal/core/domain"
	"go-hex-forum/pkg/svcerr"
)

type PostRepository interface {
	SavePost(ctx context.Context, post *domain.Post, userID int64) (int64, error)
	GetActivePosts(ctx context.Context, pagination *domain.Pagination) ([]domain.Post, error)
	GetArchivedPosts(ctx context.Context, pagination *domain.Pagination) ([]domain.Post, error)
	GetPostByID(ctx context.Context, postID int64) (domain.Post, error)
	ArchiveExpiredPosts(ctx context.Context) error
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
	return &PostService{postRepo, imageStorage}
}

func (s *PostService) CreateNewPost(ctx context.Context, post *domain.Post, imageData []byte) (int64, error) {
	const op = "PostService.CreateNewPost"
	// validation
	if post.Title == "" || post.Content == "" {
		err := fmt.Errorf("%s: title and content are not provided", op)
		return -1, svcerr.NewError("title and content are required", err, svcerr.ErrBadRequest)
	}
	post.CreatedAt = time.Now().UTC()
	// initial expiration
	post.ExpiresAt = post.CreatedAt.Add(10 * time.Minute)

	if len(imageData) > 0 {
		url, err := s.imageStorage.UploadImage(ctx, post.PostAuthor.ID, imageData)
		if err != nil {
			raw := fmt.Errorf("%s: upload image failed: %w", op, err)
			return -1, svcerr.NewError("failed to upload image", raw, svcerr.ErrInternal)
		}
		post.ImagePath = url
	}

	id, err := s.postRepo.SavePost(ctx, post, post.PostAuthor.ID)
	if err != nil {
		raw := fmt.Errorf("%s: save post failed: %w", op, err)
		return -1, svcerr.NewError("failed to save post", raw, svcerr.ErrInternal)
	}
	return id, nil
}

func (s *PostService) GetActivePosts(ctx context.Context) ([]domain.Post, error) {
	posts, err := s.postRepo.GetActivePosts(ctx, &domain.Pagination{Page: 1, PageSize: 10})
	if err != nil {
		raw := fmt.Errorf("PostService.GetActivePosts: %w", err)
		return nil, svcerr.NewError("failed to get active posts", raw, svcerr.ErrInternal)
	}
	return posts, nil
}

func (s *PostService) GetArchivedPosts(ctx context.Context) ([]domain.Post, error) {
	posts, err := s.postRepo.GetArchivedPosts(ctx, &domain.Pagination{Page: 1, PageSize: 10})
	if err != nil {
		raw := fmt.Errorf("PostService.GetArchivedPosts: %w", err)
		return nil, svcerr.NewError("failed to get archived posts", raw, svcerr.ErrInternal)
	}
	return posts, nil
}

func (s *PostService) GetPostByID(ctx context.Context, postID int64) (domain.Post, error) {
	const op = "PostService.GetPostByID"
	post, err := s.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		raw := fmt.Errorf("%s: %w", op, err)
		return domain.Post{}, svcerr.NewError("post not found", raw, svcerr.ErrNotFound)
	}
	return post, nil
}

func (s *PostService) ArchiveExpiredPostsWorker(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := s.postRepo.ArchiveExpiredPosts(ctx)
				if err != nil {
					// optionally log or handle
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
