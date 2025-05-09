package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/ports/http/httperror"
)

type PostService interface {
	CreateNewPost(ctx context.Context, post *domain.Post, imageData []byte) (int64, error)
	GetActivePosts(context.Context) ([]domain.Post, error)
	GetArchivedPosts(context.Context) ([]domain.Post, error)
	GetPostByID(ctx context.Context, postID int64) (domain.Post, error)
}

type PostHandler struct {
	postService PostService
}

func NewPostHandler(postService PostService) *PostHandler {
	return &PostHandler{postService}
}

func (h *PostHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /post", h.CreateNewPost)
}

func (h *PostHandler) CreateNewPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		httperror.WriteError(w, err)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	file, _, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		httperror.WriteError(w, err)
		return
	}
	if file != nil {
		defer file.Close()
	}

	var imageData []byte
	if file != nil {
		imageData, err = io.ReadAll(file)
		if err != nil {
			httperror.WriteError(w, err)
			return
		}
	}

	session, ok := r.Context().Value("session").(*domain.Session)
	if !ok {
		httperror.WriteError(w, err)
		return
	}

	post := &domain.Post{
		PostAuthor: session.User,
		Title:      title,
		Content:    content,
	}
	fmt.Printf("%v", post)

	id, err := h.postService.CreateNewPost(r.Context(), post, imageData)
	if err != nil {
		httperror.WriteError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}
