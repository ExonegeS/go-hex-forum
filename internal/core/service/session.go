package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
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
	GetAvatarLink() string
	GetName() string
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
	// Генерируем токен
	plainToken, err := generateSessionToken()
	if err != nil {
		return "", err
	}

	// Хэшируем SHA-256
	h := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(h[:])

	now := time.Now().UTC()
	rec := domain.Session{
		TokenHash: tokenHash,
		// User:      user,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(s.cfg.DefaultTTL)),
	}

	if err = s.sessionRepo.Store(context.Background(), rec); err != nil {
		return "", err
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
