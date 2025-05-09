package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/ports/http/httperror"
	"go-hex-forum/pkg/svcerr"
)

type FrontendHandler struct {
	templates      *template.Template
	postService    PostService
	commentService CommentService
	sessionService SessionService
	logger         *slog.Logger
}

func NewFrontendHandler(postService PostService, sessionService SessionService, commentService CommentService, tpl *template.Template, logger *slog.Logger) *FrontendHandler {
	tpl.ParseGlob(filepath.Join("web", "templates", "components", "*.html"))
	return &FrontendHandler{
		templates:      tpl,
		postService:    postService,
		sessionService: sessionService,
		commentService: commentService,
		logger:         logger,
	}
}

func (h *FrontendHandler) RegisterFrontendEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/", h.ShowIndex)
	mux.HandleFunc("/create-post", h.ShowCreatePost)
	mux.HandleFunc("/post", h.CreateNewPost)
	mux.HandleFunc("/archive", h.ShowArchive)
	mux.HandleFunc("/post/{id}", h.ShowPost)
	mux.HandleFunc("/post/{id}/comment", h.CreateNewComment)
}

func (h *FrontendHandler) ShowIndex(w http.ResponseWriter, r *http.Request) {
	const op = "FrontendHandler.ShowIndex"
	h.logger.Info("handling ShowIndex", "op", op, "method", r.Method)
	posts, err := h.postService.GetActivePosts(r.Context())
	if err != nil {
		h.logger.Warn("failed to load active posts", "op", op, "err", err)
		h.renderErrorPage(w, svcerr.NewError("failed to load posts", err, svcerr.ErrInternal))
		return
	}
	session, _ := r.Context().Value("session").(*domain.Session)
	h.renderTemplate(w, "catalog.html", map[string]interface{}{
		"Session": session,
		"Posts":   posts,
	})
	h.logger.Info("rendered ShowIndex", "op", op, "count", len(posts))
}

func (h *FrontendHandler) ShowArchive(w http.ResponseWriter, r *http.Request) {
	const op = "FrontendHandler.ShowArchive"
	h.logger.Info("handling ShowArchive", "op", op, "method", r.Method)
	posts, err := h.postService.GetArchivedPosts(r.Context())
	if err != nil {
		h.logger.Warn("failed to load archived posts", "op", op, "err", err)
		h.renderErrorPage(w, err)
		return
	}
	session, _ := r.Context().Value("session").(*domain.Session)
	h.renderTemplate(w, "archive.html", map[string]interface{}{
		"Session": session,
		"Posts":   posts,
	})
	h.logger.Info("rendered ShowArchive", "op", op, "count", len(posts))
}

func (h *FrontendHandler) ShowCreatePost(w http.ResponseWriter, r *http.Request) {
	const op = "FrontendHandler.ShowCreatePost"
	h.renderTemplate(w, "create-post.html", nil)
	h.logger.Info("rendered ShowCreatePost", "op", op)
}

func (h *FrontendHandler) CreateNewPost(w http.ResponseWriter, r *http.Request) {
	const op = "FrontendHandler.CreateNewPost"
	h.logger.Info("handling CreateNewPost", "op", op)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Warn("invalid form data", "op", op, "err", err)
		h.renderErrorPage(w, svcerr.NewError("invalid form data", err, svcerr.ErrBadRequest))
		return
	}

	title, content := r.FormValue("title"), r.FormValue("content")
	h.logger.Info("form values", "op", op, "title", title)

	file, _, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		h.logger.Warn("failed to read upload", "op", op, "err", err)
		h.renderErrorPage(w, svcerr.NewError("failed to read upload", err, svcerr.ErrInternal))
		return
	}
	var imageData []byte
	if file != nil {
		defer file.Close()
		imageData, err = io.ReadAll(file)
		if err != nil {
			h.logger.Warn("failed to read image", "op", op, "err", err)
			h.renderErrorPage(w, svcerr.NewError("failed to read image", err, svcerr.ErrInternal))
			return
		}
	}

	session, ok := r.Context().Value("session").(*domain.Session)
	if !ok {
		h.logger.Warn("no session in context", "op", op)
		h.renderErrorPage(w, svcerr.NewError("not authorized", fmt.Errorf("%s: no session", op), svcerr.ErrNotAuthorized))
		return
	}

	post := &domain.Post{PostAuthor: session.User, Title: title, Content: content}
	id, err := h.postService.CreateNewPost(r.Context(), post, imageData)
	if err != nil {
		h.logger.Warn("create post failed", "op", op, "err", err)
		h.renderErrorPage(w, err)
		return
	}

	h.logger.Info("post created", "op", op, "postID", id)
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

func (h *FrontendHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
	const op = "ShowPost"
	idStr := r.URL.Path[len("/post/"):]
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.renderErrorPage(w, svcerr.NewError("invalid post id", err, svcerr.ErrBadRequest))
		return
	}

	post, err := h.postService.GetPostByID(r.Context(), postID)
	if err != nil {
		h.renderErrorPage(w, err)
		return
	}

	comments, err := h.commentService.GetByPostID(r.Context(), postID)
	if err != nil {
		h.renderErrorPage(w, err)
		return
	}

	tpl := "post.html"
	if post.IsArchived {
		tpl = "archive-post.html"
	}
	h.renderTemplate(w, tpl, map[string]interface{}{
		"UserAvatar": post.PostAuthor.AvatarURL,
		"UserName":   post.PostAuthor.Name,
		"DataTime":   post.CreatedAt,
		"PostID":     post.ID,
		"ImagePath":  post.ImagePath,
		"Title":      post.Title,
		"Content":    post.Content,
		"Comments":   comments,
	})
}

func (h *FrontendHandler) CreateNewComment(w http.ResponseWriter, r *http.Request) {
	const op = "FrontendHandler.CreateNewComment"
	h.logger.Info("handling CreateNewComment", "op", op)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Warn("invalid comment form", "op", op, "err", err)
		h.renderErrorPage(w, svcerr.NewError("invalid comment form", err, svcerr.ErrBadRequest))
		return
	}

	content := r.FormValue("comment")
	h.logger.Info("comment content received", "op", op, "len", len(content))

	var parentID *int64
	if raw := r.FormValue("parent_comment_id"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			h.logger.Warn("invalid parent_comment_id", "op", op, "raw", raw, "err", err)
			h.renderErrorPage(w, svcerr.NewError("invalid parent comment id", err, svcerr.ErrBadRequest))
			return
		}
		parentID = &id
	}

	file, _, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		h.logger.Warn("failed to read upload file", "op", op, "err", err)
		h.renderErrorPage(w, svcerr.NewError("failed to read upload", err, svcerr.ErrInternal))
		return
	}
	var imageData []byte
	if file != nil {
		defer file.Close()
		imageData, err = io.ReadAll(file)
		if err != nil {
			h.logger.Warn("failed to read comment image", "op", op, "err", err)
			h.renderErrorPage(w, svcerr.NewError("failed to read image", err, svcerr.ErrInternal))
			return
		}
	}

	idStr := r.URL.Path[len("/post/") : len(r.URL.Path)-len("/comment")]
	postID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("invalid post id in path", "op", op, "raw", idStr, "err", err)
		h.renderErrorPage(w, svcerr.NewError("invalid post id", err, svcerr.ErrBadRequest))
		return
	}

	session, ok := r.Context().Value("session").(*domain.Session)
	if !ok {
		h.logger.Warn("no session in context", "op", op)
		h.renderErrorPage(w, svcerr.NewError("not authorized", fmt.Errorf("%s: no session", op), svcerr.ErrNotAuthorized))
		return
	}

	comment := domain.Comment{
		PostID:          postID,
		ParentCommentID: parentID,
		Content:         content,
		Author:          domain.UserData{ID: session.User.ID},
		CreatedAt:       time.Now(),
	}
	if _, err := h.commentService.SaveComment(r.Context(), &comment, imageData); err != nil {
		h.logger.Warn("save comment failed", "op", op, "err", err)
		h.renderErrorPage(w, svcerr.NewError("could not save comment", err, svcerr.ErrInternal))
		return
	}

	h.logger.Info("comment created", "op", op, "commentID", comment.ID, "postID", postID)
	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}

func (h *FrontendHandler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		httperror.WriteError(w, err)
	}
}

func (h *FrontendHandler) renderErrorPage(w http.ResponseWriter, err error) {
	apiErr := httperror.FromError(err)
	w.WriteHeader(apiErr.StatusCode)
	h.templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
		"Message":    apiErr.Message,
		"StatusCode": apiErr.StatusCode,
	})
}
