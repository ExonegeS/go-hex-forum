package service

import (
	"context"
	"fmt"
	"time"

	"go-hex-forum/internal/core/domain"
	"go-hex-forum/pkg/svcerr"
)

// Transactor abstracts transactional execution
// ImageStorage uploads images and returns public URLs

type CommentRepository interface {
	SaveComment(ctx context.Context, comment *domain.Comment) (int64, error)
	GetByPostID(ctx context.Context, postID int64) ([]*domain.Comment, error)
}

type CommentPostRepo interface {
	GetPostByID(ctx context.Context, postID int64) (domain.Post, error)
	UpdateExpiresAt(ctx context.Context, postID int64, expiresAt time.Time) error
}

type CommentService struct {
	transactor   Transactor
	commentRepo  CommentRepository
	postRepo     CommentPostRepo
	imageStorage ImageStorage
}

func NewCommentService(
	tr Transactor,
	cr CommentRepository,
	pr CommentPostRepo,
	is ImageStorage,
) *CommentService {
	return &CommentService{tr, cr, pr, is}
}

func (s *CommentService) SaveComment(ctx context.Context, comment *domain.Comment, imageData []byte) (int64, error) {
	const op = "CommentService.SaveComment"

	// Проверка поста
	post, err := s.postRepo.GetPostByID(ctx, comment.PostID)
	if err != nil {
		raw := fmt.Errorf("%s: get post: %w", op, err)
		return -1, svcerr.NewError("post not found", raw, svcerr.ErrNotFound)
	}
	if post.IsArchived {
		raw := fmt.Errorf("%s: post is archived", op)
		return -1, svcerr.NewError("post is archived, new comments are prohibited", raw, svcerr.ErrBadRequest)
	}

	// Загрузка изображения
	if len(imageData) > 0 {
		url, err := s.imageStorage.UploadImage(ctx, comment.Author.ID, imageData)
		if err != nil {
			raw := fmt.Errorf("%s: upload image: %w", op, err)
			return -1, svcerr.NewError("failed to upload comment image", raw, svcerr.ErrInternal)
		}
		comment.ImagePath = url
	}

	// Транзакция: сохранение комментария + продление expires_at
	var id int64
	err = s.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		var innerErr error
		id, innerErr = s.commentRepo.SaveComment(txCtx, comment)
		if innerErr != nil {
			raw := fmt.Errorf("%s: save comment: %w", op, innerErr)
			return svcerr.NewError("failed to save comment", raw, svcerr.ErrInternal)
		}

		if err := s.postRepo.UpdateExpiresAt(txCtx, comment.PostID, time.Now().Add(15*time.Minute)); err != nil {
			raw := fmt.Errorf("%s: update post expires: %w", op, err)
			return svcerr.NewError("failed to update post expiration", raw, svcerr.ErrInternal)
		}

		return nil
	})
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *CommentService) GetByPostID(ctx context.Context, postID int64) ([]*domain.Comment, error) {
	const op = "CommentService.GetByPostID"
	comments, err := s.commentRepo.GetByPostID(ctx, postID)
	if err != nil {
		raw := fmt.Errorf("%s: get comments: %w", op, err)
		return nil, svcerr.NewError("failed to load comments", raw, svcerr.ErrInternal)
	}
	return comments, nil
}
