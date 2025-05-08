package unit

// import (
// 	"context"
// 	"fmt"
// 	"go-hex-forum/internal/core/domain"
// 	"go-hex-forum/internal/core/service"
// 	"testing"
// )

// type mockPostRepository struct {
// 	saveFunc func(ctx context.Context, post *domain.Post, userID int64) (int64, error)
// 	getFunc  func(ctx context.Context, postID int64) (domain.Post, error)
// }

// func (m *mockPostRepository) SavePost(ctx context.Context, post *domain.Post, userID int64) (int64, error) {
// 	if m.saveFunc != nil {
// 		return m.saveFunc(ctx, post, userID)
// 	}
// 	return -2, nil
// }
// func (m *mockPostRepository) GetActivePosts(ctx context.Context, pagination *domain.Pagination) ([]domain.Post, error) {
// 	return []domain.Post{}, nil
// }
// func (m *mockPostRepository) GetPostByID(ctx context.Context, postID int64) (domain.Post, error) {
// 	if m.getFunc != nil {
// 		return m.getFunc(ctx, postID)
// 	}
// 	return domain.Post{}, nil
// }

// type mockCommentRepository struct {
// 	saveFunc func(ctx context.Context, comment *domain.Comment) (int64, error)
// 	getFunc  func(ctx context.Context, postID int64) ([]*domain.Comment, error)
// }

// func (m *mockCommentRepository) SaveComment(ctx context.Context, comment *domain.Comment) (int64, error) {
// 	if m.saveFunc != nil {
// 		return m.saveFunc(ctx, comment)
// 	}
// 	return -2, nil
// }

// func (m *mockCommentRepository) GetByPostID(ctx context.Context, postID int64) ([]*domain.Comment, error) {
// 	if m.getFunc != nil {
// 		return m.getFunc(ctx, postID)
// 	}
// 	return nil, nil
// }

// type mockImageStorage struct {
// 	uploadFunc func(ctx context.Context, userID int64, data []byte) (publicURL string, err error)
// 	getUrlFunc func(userID int64, code string) string
// }

// func (m *mockImageStorage) UploadImage(ctx context.Context, userID int64, data []byte) (publicURL string, err error) {
// 	if m.uploadFunc != nil {
// 		return m.uploadFunc(ctx, userID, data)
// 	}
// 	return "", nil
// }
// func (m *mockImageStorage) GetImageURL(userID int64, code string) string {
// 	if m.getUrlFunc != nil {
// 		return m.getUrlFunc(userID, code)
// 	}
// 	return ""
// }
// func TestCreateComment_Success(t *testing.T) {
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
// 			return "http://localhost:6969/user-104/dBAvuR", nil
// 		},
// 	}

// 	service := service.NewCommentService(repoComment, repoPost, imageMock)
// 	id, err := service.SaveComment(context.Background(), &domain.Comment{
// 		PostID: 1,
// 	}, []byte(""))
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}
// 	if id != 1 {
// 		t.Fatalf("expected id 1, got %d", id)
// 	}
// }

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

// 	service := service.NewCommentService(repoComment, repoPost, imageMock)
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

// 	service := service.NewCommentService(repoComment, repoPost, imageMock)
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

// 	service := service.NewCommentService(repoComment, repoPost, imageMock)
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

// 	service := service.NewCommentService(repoComment, repoPost, imageMock)
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
