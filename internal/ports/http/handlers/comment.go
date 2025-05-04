package handlers

import (
	"encoding/json"
	"errors"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/ports/dto"
	"go-hex-forum/internal/utils"
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
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	postID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid post ID"))
		return
	}

	session, ok := r.Context().Value("session").(*domain.Session)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	comment := domain.Comment{
		PostID:    postID,
		Content:   req.Content,
		ImagePath: req.ImagePath,
		Author: domain.UserData{
			ID: session.User.ID,
		},
	}

	_, err = h.CommentService.SaveNewComment(&comment, postID, session.User.ID)
}
