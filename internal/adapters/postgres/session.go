package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/core/service"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db}
}

func (r *SessionRepository) Store(ctx context.Context, session domain.Session) error {
	const op = "SessionRepository.Store"

	return runInTx(r.db, func(tx *sql.Tx) error {
		userStmt, err := tx.PrepareContext(ctx, `
			INSERT INTO users(name, avatar_url)
			VALUES($1, $2)
			RETURNING id
		`)
		if err != nil {
			return fmt.Errorf("%s: prepare user insert: %w", op, err)
		}
		defer userStmt.Close()

		var userID int64
		err = userStmt.QueryRowContext(ctx, session.User.Name, session.User.AvatarURL).Scan(&userID)
		if err != nil {
			return fmt.Errorf("%s: user insert: %w", op, err)
		}

		sessionStmt, err := tx.PrepareContext(ctx, `
			INSERT INTO sessions(
				session_hash, 
				user_id, 
				created_at, 
				expires_at
			) VALUES($1, $2, $3, $4)
		`)
		if err != nil {
			return fmt.Errorf("%s: prepare session insert: %w", op, err)
		}
		defer sessionStmt.Close()

		_, err = sessionStmt.ExecContext(ctx,
			session.TokenHash,
			userID,
			session.CreatedAt,
			session.ExpiresAt,
		)
		if err != nil {
			return fmt.Errorf("%s: session insert: %w", op, err)
		}

		return nil
	})
}

func (r *SessionRepository) GetByHashedToken(ctx context.Context, hashedToken string) (*domain.Session, error) {
	const op = "SessionRepository.GetByHashedToken"

	query := `
		SELECT 
			s.id, 
			s.session_hash, 
			s.created_at, 
			s.expires_at,
			u.id AS user_id,
			u.name,
			u.avatar_url
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.session_hash = $1
	`

	var session domain.Session
	err := r.db.QueryRowContext(ctx, query, hashedToken).Scan(
		&session.ID,
		&session.TokenHash,
		&session.CreatedAt,
		&session.ExpiresAt,
		&session.User.ID,
		&session.User.Name,
		&session.User.AvatarURL,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%s: session get: %w", op, err)
	}

	return &session, nil
}

func (r *SessionRepository) UpdateByToken(ctx context.Context, hashedToken string, updateFn func(*domain.Session) (bool, error)) error {
	const op = "SessionRepository.UpdateByToken"

	return runInTx(r.db, func(tx *sql.Tx) error {
		query := `SELECT
			s.id,
			s.session_hash,
			s.expires_at,
			u.id AS user_id,
			u.name,
			u.avatar_url
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		WHERE s.session_hash = $1 FOR UPDATE`

		var (
			id            int64
			tokenHash     string
			expiresAt     time.Time
			userId        int64
			userName      string
			userAvatarURL string
		)
		err := tx.QueryRowContext(ctx, query, hashedToken).Scan(&id, &tokenHash, &expiresAt, &userId, &userName, &userAvatarURL)
		if err != nil {
			return fmt.Errorf("%s: session get: %w", op, err)
		}

		session := &domain.Session{
			ID:        id,
			TokenHash: tokenHash,
			ExpiresAt: expiresAt,
			User: domain.UserData{
				ID:        userId,
				Name:      userName,
				AvatarURL: userAvatarURL,
			},
		}

		updated, err := updateFn(session)
		if err != nil {
			return fmt.Errorf("%s: updateFn: %w", op, err)
		}

		if !updated {
			return nil
		}

		_, err = tx.ExecContext(ctx,
			"UPDATE sessions SET expires_at = $1 WHERE session_hash = $2",
			session.ExpiresAt, session.TokenHash)
		if err != nil {
			return fmt.Errorf("%s: update session: %w", op, err)
		}

		_, err = tx.ExecContext(ctx,
			"UPDATE users SET name = $1, avatar_url = $2 WHERE id = $3",
			session.User.Name, session.User.AvatarURL, session.User.ID)
		if err != nil {
			return fmt.Errorf("%s: update user: %w", op, err)
		}

		return nil
	})
}

func runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	const op = "SessionRepository.runInTx"

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	err = fn(tx)
	if err == nil {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("%s: commit transaction: %w", op, err)
		}
		return nil
	}

	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return errors.Join(fmt.Errorf("%s: fn failed: %w", op, err), fmt.Errorf("%s: rollback failed: %w", op, rollbackErr))
	}

	return fmt.Errorf("%s: fn failed: %w", op, err)
}
