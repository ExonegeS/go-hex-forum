package dto

import (
	"go-hex-forum/internal/core/domain"
	"time"
)

type CreateCommentRequest struct {
	Content   string `json:"content"`
	ImagePath string `json:"image_path,omitempty"`
}

type CommentResponse struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Author    UserData  `json:"author"`
}

type UserData struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func RequestToDomain(req CreateCommentRequest, postID, userID int64) domain.Comment {
	return domain.Comment{
		PostID:    postID,
		Content:   req.Content,
		ImagePath: req.ImagePath,
		Author: domain.UserData{
			ID: userID,
		},
	}
}
