package domain

import "time"

type Session struct {
	ID        int64
	TokenHash string // Захэшированный токен для безопасности
	User      UserData
	ExpiresAt time.Time
	CreatedAt time.Time
}

type UserData struct {
	ID        int64
	Name      string
	AvatarURL string
}

func (s Session) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}
