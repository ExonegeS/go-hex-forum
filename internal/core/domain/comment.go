package domain

import "time"

type Comment struct {
	ID        int64
	PostID    int64
	User      UserData
	Content   string
	ImagePath string // S3 path для изображений комментария
	CreatedAt time.Time
	Author    UserData
}
