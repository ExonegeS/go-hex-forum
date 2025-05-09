package dto

import (
	"time"

	"go-hex-forum/internal/core/domain"
)

type PostResponse struct {
	ID              int64     `json:"id"`
	AuthorName      string    `json:"author_name"`
	AuthorAvatarURL string    `json:"author_avatar_url"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	ImageURL        string    `json:"image_url"`
	CreatedAt       time.Time `json:"created_at"`
	IsArchived      bool      `json:"is_archived"`
}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	// файл картинки обычно передают отдельно, либо URL внешнего загрузчика
}

func ToPostResponse(p *domain.Post) *PostResponse {
	return &PostResponse{
		ID:              p.ID,
		AuthorName:      p.PostAuthor.Name,
		AuthorAvatarURL: p.PostAuthor.AvatarURL,
		Title:           p.Title,
		Content:         p.Content,
		ImageURL:        p.ImagePath, // при необходимости конверсия S3 path → публичный URL
		CreatedAt:       p.CreatedAt,
		IsArchived:      p.IsArchived,
	}
}

func FromCreatePostRequest(req *CreatePostRequest, author domain.UserData) *domain.Post {
	return &domain.Post{
		PostAuthor: author,
		Title:      req.Title,
		Content:    req.Content,
		// ImagePath заполняется позднее — например, после загрузки в S3
	}
}
