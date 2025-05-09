package service

import (
	"context"
	"testing"
	"time"

	"go-hex-forum/config"
	"go-hex-forum/internal/core/domain"
)

type mockSessionRepository struct {
	saveFunc func(context.Context, domain.Session) error
	// GetByHashedToken(context.Context, string) (*domain.Session, error)
	// UpdateByToken(context.Context, string, func(*domain.Session) (bool, error)) error
}

func (m *mockSessionRepository) GetByHashedToken(context.Context, string) (*domain.Session, error) {
	return &domain.Session{
		ID:        1,
		TokenHash: "",
		User: domain.UserData{
			ID:        0,
			Name:      "Ex",
			AvatarURL: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		},
	}, nil
}

func (m *mockSessionRepository) Store(ctx context.Context, s domain.Session) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, s)
	}
	return nil
}

func (m *mockSessionRepository) UpdateByToken(context.Context, string, func(*domain.Session) (bool, error)) error {
	return nil
}

type mockUserDataProvider struct {
	getFn func(ttl time.Duration) (domain.UserData, error)
}

func (m *mockUserDataProvider) GetUserData(ttl time.Duration) (domain.UserData, error) {
	if m.getFn != nil {
		return m.getFn(ttl)
	}
	return domain.UserData{
		ID:        0,
		Name:      "Ex",
		AvatarURL: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
	}, nil
}

func TestStoreNewSession_Success(t *testing.T) {
	sessionRepo := &mockSessionRepository{
		saveFunc: func(ctx context.Context, s domain.Session) error {
			return nil
		},
	}
	userDataProvider := &mockUserDataProvider{}
	mockConfig := config.SessionConfig{
		DefaultTTL:    1 * time.Minute,
		MaxNameLength: 10,
	}

	service := NewSessionService(sessionRepo, time.Now, userDataProvider, mockConfig)
	if service == nil {
		t.Fatalf("NewSessionService returned nil")
	}
	_, err := service.StoreNewSession(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
