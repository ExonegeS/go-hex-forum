package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"go-hex-forum/config"
	"go-hex-forum/internal/core/domain"
	"os"
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

// func (s *SessionsService) CreateSession(ctx context.Context, ip, userAgent string) (string, error) {

// 	session := entity.Session{
// 		// Name:        s.userDataAPI.GetName(),
// 		// Avatar:      s.userDataAPI.GetAvatarLink(),
// 		ExpiresAt:   s.timeSource().Add(s.cfg.DefaultTTL),
// 		Fingerprint: generateFingerprint(ip, userAgent),
// 		IP:          ip,
// 		UserAgent:   userAgent,
// 	}

// 	sessionID, err := s.sessionsRepo.Store(ctx, session)
// 	if err != nil {
// 		return "", fmt.Errorf("session storage failed: %w", err)
// 	}

// 	return sessionToken, nil
// }

func (s *SessionService) StoreNewSession() (plainToken string, err error) {
	// 1. Генерируем raw-токен
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", err
	}
	plainToken = hex.EncodeToString(raw)

	// 2. Хэшируем SHA-256
	h := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(h[:])

	now := time.Now().UTC()
	rec := domain.Session{
		TokenHash: tokenHash,
		// User:      user,
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(s.cfg.DefaultTTL) * time.Second),
	}

	// 3. Сохраняем в хранилище
	if err = s.sessionRepo.Store(context.Background(), rec); err != nil {
		return "", err
	}

	return plainToken, nil
}

// ValidateSession по plain-token возвращает Session или ошибку
func (s *SessionService) ValidateSession(plainToken string) (*domain.Session, error) {
	// 1. Хэшируем вход
	h := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(h[:])

	// 2. Загружаем из репозитория
	rec, err := s.sessionRepo.GetByHashedToken(context.Background(), tokenHash)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// 3. Проверяем истечение
	if time.Now().UTC().After(rec.ExpiresAt) {
		// можно и удалить сразу: s.repo.DeleteByHash(tokenHash)
		return nil, ErrSessionExpired
	}

	return &domain.Session{
		ID:        rec.ID,
		TokenHash: rec.TokenHash,
		User:      rec.User,
		CreatedAt: rec.CreatedAt,
		ExpiresAt: rec.ExpiresAt,
	}, nil
}

func generateFingerprint(ip, userAgent string) string {
	secret := []byte(os.Getenv("FINGERPRINT_PEPPER"))
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(ip))
	h.Write([]byte(userAgent))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
