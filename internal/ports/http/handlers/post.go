package handlers

import (
	"context"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/utils"
	"io"
	"net/http"
)

type PostService interface {
	CreateNewPost(ctx context.Context, post *domain.Post, imageData []byte) (int64, error)
	GetActivePosts(context.Context) ([]domain.Post, error)
	GetPostByID(ctx context.Context, postID int64) (domain.Post, error)
	UploadImage(ctx context.Context, userID int64, imageData []byte) (string, error)
}

type PostHandler struct {
	postService PostService
}

func NewPostHandler(postService PostService) *PostHandler {
	return &PostHandler{postService}
}

func (h *PostHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /post", h.CreateNewPost)
	mux.HandleFunc("POST /image", h.UploadImage)
}

func (h *PostHandler) CreateNewPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse form: %w", err))
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	file, _, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error reading image: %w", err))
		return
	}
	if file != nil {
		defer file.Close()
	}

	var imageData []byte
	if file != nil {
		imageData, err = io.ReadAll(file)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error reading image: %w", err))
			return
		}
	}

	session, ok := r.Context().Value("session").(*domain.Session)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated", session))
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
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

func (h *PostHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error reading image data: %w", err))
		return
	}

	session, ok := r.Context().Value("session").(*domain.Session)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	url, err := h.postService.UploadImage(r.Context(), session.User.ID, imageData)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"url": url,
	})
}
