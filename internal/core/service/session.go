package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"go-hex-forum/config"
	"go-hex-forum/internal/core/domain"
	"time"
)

type SessionRepository interface {
	Store(context.Context, domain.Session) error
	GetByHashedToken(context.Context, string) (*domain.Session, error)
	// DeleteExpired(context.Context, time.Time) (int64, error)
	// UpdateByID(context.Context, int64, func(*domain.Session) (bool, error)) error
}

type UserDataProvider interface {
	GetUserData(ttl time.Duration) (domain.UserData, error)
}

type SessionService struct {
	sessionRepo SessionRepository
	timeSource  func() time.Time
	userDataAPI UserDataProvider
	cfg         config.SessionConfig
}

func NewSessionService(sessionsRepo SessionRepository, timeSource func() time.Time, userDataAPI UserDataProvider, cfg config.SessionConfig) *SessionService {
	return &SessionService{
		sessionsRepo,
		timeSource,
		userDataAPI,
		cfg,
	}
}

func (s *SessionService) StoreNewSession() (string, error) {
	// Get user data with session TTL
	userData, err := s.userDataAPI.GetUserData(s.cfg.DefaultTTL)
	if err != nil {
		return "", fmt.Errorf("failed to get user data: %w", err)
	}

	fmt.Println(userData)

	// Generate session token
	plainToken, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	// Create session record
	h := sha256.Sum256([]byte(plainToken))
	session := domain.Session{
		TokenHash: hex.EncodeToString(h[:]),
		User:      userData,
		CreatedAt: s.timeSource(),
		ExpiresAt: s.timeSource().Add(s.cfg.DefaultTTL),
	}

	if err := s.sessionRepo.Store(context.Background(), session); err != nil {
		return "", fmt.Errorf("failed to store session: %w", err)
	}

	return plainToken, nil
}

func (s *SessionService) GetSessionByToken(plainToken string) (*domain.Session, error) {
	// Хэшируем токен
	h := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(h[:])

	// Загружаем сессию из репозитория
	session, err := s.sessionRepo.GetByHashedToken(context.Background(), tokenHash)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Проверяем истечение
	if session.IsExpired() {
		// можно и удалить сразу: s.repo.DeleteByHash(tokenHash)
		return nil, ErrSessionExpired
	}

	return session, nil
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32) // 256 бит
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
