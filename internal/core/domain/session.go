package domain

import "time"

type Session struct {
	ID        int64
	TokenHash string // Захэшированный токен для безопасности
	User      UserData
	Name      string // Имя из API (может быть изменено пользователем)
	ExpiresAt time.Time
	CreatedAt time.Time
}

type UserData struct {
	Name      string
	AvatarURL string
}

func (s Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
