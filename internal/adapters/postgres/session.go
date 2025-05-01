package postgres

import (
	"context"
	"database/sql"
	"go-hex-forum/internal/core/domain"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db}
}

func (r *SessionRepository) Store(ctx context.Context, session domain.Session) error {
	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO sessions (session_hash,avatar_url,username,created_at,expires_at)
	VALUES ($1,$2,$3,$4,$5)`)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		session.TokenHash,
		session.User.AvatarURL,
		session.User.Name,
		session.CreatedAt,
		session.ExpiresAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) GetByHashedToken(ctx context.Context, hashedToken string) (*domain.Session, error) {
	var session domain.Session

	stmt, err := r.db.PrepareContext(ctx, `SELECT id,session_hash,avatar_url,username,created_at,expires_at FROM sessions
	WHERE session_hash=$1`)
	if err != nil {
		return &session, err
	}

	stmt.QueryRowContext(ctx, hashedToken).Scan(
		&session.ID,
		&session.TokenHash,
		&session.User.AvatarURL,
		&session.User.Name,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	return &session, nil
}
