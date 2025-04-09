package unit

import (
	"context"
	"errors"
	"testing"

	"github.com/ExonegeS/go-hex-forum/internal/application"
	"github.com/ExonegeS/go-hex-forum/internal/domain/model"
)

type mockXRepository struct {
	saveFunc func(ctx context.Context, x model.X) (int, error)
}

func (m *mockXRepository) Save(ctx context.Context, x model.X) (int, error) {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, x)
	}
	return 0, nil
}

func TestCreateX_Success(t *testing.T) {
	repo := &mockXRepository{
		saveFunc: func(ctx context.Context, x model.X) (int, error) {
			return 1, nil
		},
	}
	service := application.NewXService(repo)
	id, err := service.CreateX(context.Background(), "test data")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 1 {
		t.Fatalf("expected id 1, got %d", id)
	}
}

func TestCreateX_Failure(t *testing.T) {
	repo := &mockXRepository{
		saveFunc: func(ctx context.Context, x model.X) (int, error) {
			return -1, errors.New("repository error")
		},
	}
	service := application.NewXService(repo)
	id, err := service.CreateX(context.Background(), "test data")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
	if id != -1 {
		t.Fatalf("expected id -1, got %d", id)
	}
}
