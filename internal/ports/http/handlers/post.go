package handlers

import (
	"context"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/utils"
	"net/http"
)

type PostService interface {
	CreateNewPost(ctx context.Context, post *domain.Post) (int64, error)
}

type PostHandler struct {
	postService PostService

	// templates *template.Template
}

func NewPostHandler(postService PostService) *PostHandler {
	return &PostHandler{postService}
}

func (h *PostHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /comment", h.CreateNewPost)
}

func (h *PostHandler) CreateNewPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // max ~10MB
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse multipart form: %w", err))
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	// file, fileHeader, err := r.FormFile("image")
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("image file is required: %w", err))
	// 	return
	// }
	// defer file.Close()

	if title == "" || content == "" {
		utils.WriteMessage(w, http.StatusBadRequest, "title and content are required")
		return
	}

	post := &domain.Post{
		Title:   title,
		Content: content,
		// ImagePath: imagePath,
	}

	id, err := h.postService.CreateNewPost(r.Context(), post)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}
