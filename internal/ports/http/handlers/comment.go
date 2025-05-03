package handlers

import (
	"encoding/json"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/ports/dto"
	"net/http"
	"strconv"
)

type CommentService interface {
	SaveNewComment(comment *domain.Comment, postID int64, userID int64) (int64, error)
}

type CommentHandler struct {
	CommentService CommentService
}

func NewCommentHandler(service CommentService) *CommentHandler {
	return &CommentHandler{service}
}

func (h *CommentHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /post/{id}/comment", h.CreateNewComment)
}

func (h *CommentHandler) CreateNewComment(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	postID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value("user_id").(int64)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	comment := domain.Comment{
		PostID:    postID,
		Content:   req.Content,
		ImagePath: req.ImagePath,
		Author: domain.UserData{
			ID: userID,
		},
	}

	_, err = h.CommentService.SaveNewComment(&comment, postID, userID)
}
