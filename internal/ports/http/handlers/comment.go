package handlers

import (
	"context"
	"errors"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/utils"
	"io"
	"net/http"
	"strconv"
	"time"
)

type CommentService interface {
	SaveComment(ctx context.Context, comment *domain.Comment, imageData []byte) (int64, error)
	GetByPostID(ctx context.Context, postID int64) ([]*domain.Comment, error)
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
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("could not parse form"))
	}

	content := r.FormValue("comment")

	var parentCommentID *int64
	if rawParentID := r.FormValue("parent_comment_id"); rawParentID != "" {
		id, err := strconv.ParseInt(rawParentID, 10, 64)
		if err == nil {
			parentCommentID = &id
		}
	}

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
		PostID:          postID,
		ParentCommentID: parentCommentID,
		Content:         content,
		Author: domain.UserData{
			ID: session.User.ID,
		},
		CreatedAt: time.Now(),
	}

	_, err = h.CommentService.SaveComment(r.Context(), &comment, imageData)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}
