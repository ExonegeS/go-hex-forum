package handlers

import (
	"context"
	"fmt"
	"go-hex-forum/internal/utils"
	"io"
	"net/http"
)

type PostService interface {
	CreateNewPost(ctx context.Context, title, content string, imageData []byte, userID int64) (int64, error)
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
	defer file.Close()

	var imageData []byte
	if file != nil {
		imageData, err = io.ReadAll(file)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error reading image data: %w", err))
			return
		}
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	id, err := h.postService.CreateNewPost(r.Context(), title, content, imageData, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

func (h *PostHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	imageData, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error reading image data: %w", err))
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("user not authenticated"))
		return
	}

	url, err := h.postService.UploadImage(r.Context(), userID, imageData)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"url": url,
	})
}
