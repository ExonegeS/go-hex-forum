package domain

import "time"

type Comment struct {
	ID        int64
	PostID    int64
	ParentID  *int64 // ID родительского комментария
	Content   string
	ImagePath string // S3 path для изображений комментария
	CreatedAt time.Time
	Author    UserData
}
