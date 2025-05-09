package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"go-hex-forum/config"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/pkg/svcerr"
)

type SessionRepository interface {
	Store(context.Context, domain.Session) error
	GetByHashedToken(context.Context, string) (*domain.Session, error)
	UpdateByToken(context.Context, string, func(*domain.Session) (bool, error)) error
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

func (s *SessionService) StoreNewSession(ctx context.Context) (string, error) {
	const op = "SessionService.StoreNewSession"

	userData, err := s.userDataAPI.GetUserData(s.cfg.DefaultTTL)
	if err != nil {
		return "", svcerr.NewError("failed to get user data", fmt.Errorf("%s: %w", op, err), svcerr.ErrInternal)
	}

	plainToken, err := generateSessionToken()
	if err != nil {
		return "", svcerr.NewError("failed to generate token", fmt.Errorf("%s: %w", op, err), svcerr.ErrInternal)
	}

	h := sha256.Sum256([]byte(plainToken))
	session := domain.Session{
		TokenHash: hex.EncodeToString(h[:]),
		User:      userData,
		CreatedAt: s.timeSource(),
		ExpiresAt: s.timeSource().Add(s.cfg.DefaultTTL),
	}

	if err := s.sessionRepo.Store(ctx, session); err != nil {
		return "", svcerr.NewError("failed to store session", fmt.Errorf("%s: %w", op, err), svcerr.ErrInternal)
	}

	return plainToken, nil
}

func (s *SessionService) GetSessionByToken(ctx context.Context, plainToken string) (*domain.Session, error) {
	const op = "SessionService.GetSessionByToken"

	h := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(h[:])

	session, err := s.sessionRepo.GetByHashedToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return nil, svcerr.NewError("session not found", ErrSessionNotFound, svcerr.ErrBadRequest)
		}
		return nil, svcerr.NewError("session not found", fmt.Errorf("%s: %w", op, err), svcerr.ErrInternal)
	}

	if session.IsExpired() {
		return nil, svcerr.NewError("session expired", ErrSessionExpired, svcerr.ErrBadRequest)
	}

	return session, nil
}

func (s *SessionService) UpdateUserName(ctx context.Context, plainToken string, username string) error {
	const op = "SessionService.UpdateUserName"

	if len(username) > 16 {
		return svcerr.NewError("too long name", errors.New("too long name maximux 16 characters"), svcerr.ErrBadRequest)
	} else if len(username) < 3 {
		return svcerr.NewError("too long name", errors.New("too short name min 3 characters"), svcerr.ErrBadRequest)
	}

	h := sha256.Sum256([]byte(plainToken))
	tokenHash := hex.EncodeToString(h[:])

	err := s.sessionRepo.UpdateByToken(ctx, tokenHash, func(s *domain.Session) (bool, error) {
		if username != "" {
			s.User.Name = username
			return true, nil
		}
		return false, fmt.Errorf("no fields were updated")
	})
	if err != nil {
		return svcerr.NewError("update failed", fmt.Errorf("%s: %w", op, err), svcerr.ErrBadRequest)
	}

	return nil
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
