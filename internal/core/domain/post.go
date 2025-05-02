package domain

import "time"

type Post struct {
	ID         int64
	PostAuthor UserData
	Title      string
	Content    string
	ImagePath  string // S3 object path (пример: "posts/abc123.jpg")
	CreatedAt  time.Time
	ExpiresAt  time.Time
	IsArchived bool
}

func (p *Post) IsExpired() bool {
	// logic here
	return false
}

func (p *Post) ArchivePostIfExpired() {
	if p.IsExpired() {
		p.IsArchived = true
	}
}
