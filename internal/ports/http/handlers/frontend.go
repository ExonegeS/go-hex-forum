package handlers

import (
	"errors"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/utils"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type FrontendHandler struct {
	templates      *template.Template
	postService    PostService
	commentService CommentService
	sessionService SessionService
}

func NewFrontendHandler(postService PostService, sessionService SessionService, commentService CommentService, tpl *template.Template) *FrontendHandler {
	// Parse component templates
	tpl.ParseGlob(filepath.Join("web", "templates", "components", "*.html"))
	// tpl.ParseGlob(filepath.Join("templates", "posts", "*.html"))

	return &FrontendHandler{
		templates:      tpl,
		postService:    postService,
		sessionService: sessionService,
		commentService: commentService,
	}
}

func (h *FrontendHandler) RegisterFrontendEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.ShowIndex)
	mux.HandleFunc("GET /create-post", h.CreatePost)
	mux.HandleFunc("GET /post/{id}", h.ShowPost)
	// mux.HandleFunc("POST /post", h.CreatePost)
	// mux.HandleFunc("GET /upload", h.ShowUploadForm)
	mux.HandleFunc("GET /archive", h.ShowArchive)
}

func (h *FrontendHandler) ShowIndex(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetActivePosts(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	// Get session from context
	session, _ := r.Context().Value("session").(*domain.Session)

	fmt.Println(session)

	h.renderTemplate(w, "catalog.html", struct {
		Session *domain.Session
		Posts   []domain.Post
	}{
		Session: session,
		Posts:   posts,
	})
}

func (h *FrontendHandler) ShowArchive(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetArchivedPosts(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	// Get session from context
	session, _ := r.Context().Value("session").(*domain.Session)

	fmt.Println(session)

	h.renderTemplate(w, "archive.html", struct {
		Session *domain.Session
		Posts   []domain.Post
	}{
		Session: session,
		Posts:   posts,
	})
}

func (h *FrontendHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Get session from context
	session, _ := r.Context().Value("session").(*domain.Session)

	fmt.Println(session)

	h.renderTemplate(w, "create-post.html", struct {
	}{})
}

func (h *FrontendHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("invalid post ID"))
		return
	}

	post, err := h.postService.GetPostByID(r.Context(), postID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("post was not found", err.Error()))
		return
	}

	comments, err := CommentService.GetByPostID(h.commentService, r.Context(), postID)
	if err != nil {
		log.Print(err)
	}
	fmt.Print("HERERR", post.CreatedAt)
	h.renderTemplate(w, "post.html", struct {
		UserAvatar string
		UserName   string
		DataTime   time.Time
		PostID     int64
		ImagePath  string
		Title      string
		Content    string
		Comments   []*domain.Comment
	}{
		UserAvatar: post.PostAuthor.AvatarURL,
		UserName:   post.PostAuthor.Name,
		DataTime:   post.CreatedAt,
		PostID:     post.ID,
		ImagePath:  post.ImagePath,
		Title:      post.Title,
		Content:    post.Content,
		Comments:   comments,
	})
}

func (h *FrontendHandler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}
