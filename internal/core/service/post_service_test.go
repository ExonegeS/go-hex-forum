package service

import (
	"context"
	"fmt"
	"testing"

	"go-hex-forum/internal/core/domain"
)

func TestCreatePost_Success(t *testing.T) {
	repoPost := &mockPostRepository{
		saveFunc: func(ctx context.Context, post *domain.Post, userID int64) (int64, error) {
			return 1, nil
		},
		getFunc: func(ctx context.Context, postID int64) (domain.Post, error) {
			return domain.Post{ID: 1}, nil
		},
		archieveExpiredfunc: func(ctx context.Context) error {
			return nil
		},
	}
	imageMock := &mockImageStorage{
		uploadFunc: func(ctx context.Context, userID int64, data []byte) (publicURL string, err error) {
			return "http://localhost:6969/user-104/dBAvuR", nil
		},
	}

	service := NewPostService(repoPost, imageMock)
	id, err := service.CreateNewPost(context.Background(), &domain.Post{
		Title:   "title",
		Content: "content",
	}, []byte(""))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 1 {
		t.Fatalf("expected id 1, got %d", id)
	}
}

func TestCreatePost_Fail(t *testing.T) {
	repoPost := &mockPostRepository{
		saveFunc: func(ctx context.Context, post *domain.Post, userID int64) (int64, error) {
			return 1, nil
		},
		getFunc: func(ctx context.Context, postID int64) (domain.Post, error) {
			return domain.Post{ID: 1}, nil
		},
	}
	imageMock := &mockImageStorage{
		uploadFunc: func(ctx context.Context, userID int64, data []byte) (publicURL string, err error) {
			return "http://localhost:6969/user-104/dBAvuR", nil
		},
	}

	service := NewPostService(repoPost, imageMock)
	id, err := service.CreateNewPost(context.Background(), &domain.Post{
		Title:   "",
		Content: "",
	}, []byte(""))
	if err == nil {
		t.Fatalf("expected error, got none")
	}
	if id != -1 {
		t.Fatalf("expected id -1, got %d", id)
	}
}

func TestCreatePost_Fail_postRepo(t *testing.T) {
	repoPost := &mockPostRepository{
		saveFunc: func(ctx context.Context, post *domain.Post, userID int64) (int64, error) {
			return -1, fmt.Errorf("error")
		},
		getFunc: func(ctx context.Context, postID int64) (domain.Post, error) {
			return domain.Post{ID: 1}, nil
		},
	}
	imageMock := &mockImageStorage{
		uploadFunc: func(ctx context.Context, userID int64, data []byte) (publicURL string, err error) {
			return "http://localhost:6969/user-104/dBAvuR", nil
		},
	}

	service := NewPostService(repoPost, imageMock)
	id, err := service.CreateNewPost(context.Background(), &domain.Post{
		Title:   "title",
		Content: "content",
	}, []byte(""))
	if err == nil {
		t.Fatalf("expected error, got none")
	}
	if id != -1 {
		t.Fatalf("expected id -1, got %d", id)
	}
}

func TestCreatePost_Fail_imgSource(t *testing.T) {
	repoPost := &mockPostRepository{
		saveFunc: func(ctx context.Context, post *domain.Post, userID int64) (int64, error) {
			return 1, nil
		},
		getFunc: func(ctx context.Context, postID int64) (domain.Post, error) {
			return domain.Post{ID: 1}, nil
		},
	}
	imageMock := &mockImageStorage{
		uploadFunc: func(ctx context.Context, userID int64, data []byte) (publicURL string, err error) {
			return "", fmt.Errorf("error")
		},
	}

	service := NewPostService(repoPost, imageMock)
	id, err := service.CreateNewPost(context.Background(), &domain.Post{
		Title:   "title",
		Content: "content",
	}, []byte("X"))
	if err == nil {
		t.Fatalf("expected error, got none")
	}
	if id != -1 {
		t.Fatalf("expected id -1, got %d", id)
	}
}

// func TestCreateComment_Fail_postRepo(t *testing.T) {
// 	repoPost := &mockPostRepository{
// 		getFunc: func(ctx context.Context, postID int64) (domain.Post, error) {
// 			return domain.Post{}, fmt.Errorf("error")
// 		},
// 	}
// 	repoComment := &mockCommentRepository{
// 		saveFunc: func(ctx context.Context, comment *domain.Comment) (int64, error) {
// 			return 1, nil
// 		},
// 	}
// 	imageMock := &mockImageStorage{
// 		uploadFunc: func(ctx context.Context, userID int64, data []byte) (publicURL string, err error) {
// 			return "http://localhost:6969/user-104/dBAvuR", nil
// 		},
// 	}

// 	service := NewCommentService(repoComment, repoPost, imageMock)
// 	id, err := service.SaveComment(context.Background(), &domain.Comment{
// 		PostID: 1,
// 	}, []byte(""))
// 	if err == nil {
// 		t.Fatalf("expected error, got none")
// 	}
// 	if id != -1 {
// 		t.Fatalf("expected id -1, got %d", id)
// 	}
// }

// func TestCreateComment_Fail_postComm(t *testing.T) {
// 	repoPost := &mockPostRepository{
// 		getFunc: func(ctx context.Context, postID int64) (domain.Post, error) {
// 			return domain.Post{ID: 1}, nil
// 		},
// 	}
// 	repoComment := &mockCommentRepository{
// 		saveFunc: func(ctx context.Context, comment *domain.Comment) (int64, error) {
// 			return -1, fmt.Errorf("error")
// 		},
// 	}
// 	imageMock := &mockImageStorage{
// 		uploadFunc: func(ctx context.Context, userID int64, data []byte) (publicURL string, err error) {
// 			return "http://localhost:6969/user-104/dBAvuR", nil
// 		},
// 	}

// 	service := NewCommentService(repoComment, repoPost, imageMock)
// 	id, err := service.SaveComment(context.Background(), &domain.Comment{
// 		PostID: 1,
// 	}, []byte(""))
// 	if err == nil {
// 		t.Fatalf("expected error, got none")
// 	}
// 	if id != -1 {
// 		t.Fatalf("expected id -1, got %d", id)
// 	}
// }

// func TestCreateComment_Success_imgSource(t *testing.T) {
// 	repoPost := &mockPostRepository{
// 		getFunc: func(ctx context.Context, postID int64) (domain.Post, error) {
// 			return domain.Post{ID: 1}, nil
// 		},
// 	}
// 	repoComment := &mockCommentRepository{
// 		saveFunc: func(ctx context.Context, comment *domain.Comment) (int64, error) {
// 			return 1, nil
// 		},
// 	}
// 	imageMock := &mockImageStorage{
// 		uploadFunc: func(ctx context.Context, userID int64, data []byte) (publicURL string, err error) {
// 			return "", fmt.Errorf("error")
// 		},
// 	}

// 	service := NewCommentService(repoComment, repoPost, imageMock)
// 	id, err := service.SaveComment(context.Background(), &domain.Comment{
// 		PostID: 1,
// 	}, []byte(""))
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}
// 	if id != 1 {
// 		t.Fatalf("expected id -1, got %d", id)
// 	}
// }

// func TestCreateComment_Fail_imgSource(t *testing.T) {
// 	repoPost := &mockPostRepository{
// 		getFunc: func(ctx context.Context, postID int64) (domain.Post, error) {
// 			return domain.Post{ID: 1}, nil
// 		},
// 	}
// 	repoComment := &mockCommentRepository{
// 		saveFunc: func(ctx context.Context, comment *domain.Comment) (int64, error) {
// 			return 1, nil
// 		},
// 	}
// 	imageMock := &mockImageStorage{
// 		uploadFunc: func(ctx context.Context, userID int64, data []byte) (publicURL string, err error) {
// 			return "", fmt.Errorf("error")
// 		},
// 	}

// 	service := NewCommentService(repoComment, repoPost, imageMock)
// 	id, err := service.SaveComment(context.Background(), &domain.Comment{
// 		PostID: 1,
// 	}, []byte("X"))
// 	if err == nil {
// 		t.Fatalf("expected error, got none")
// 	}
// 	if id != -1 {
// 		t.Fatalf("expected id -1, got %d", id)
// 	}
// }
