package handlers

import (
	"fmt"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/ports/http/httperror"
	"go-hex-forum/internal/utils"
	"html/template"
	"io"
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
	tpl.ParseGlob(filepath.Join("web", "templates", "components", "*.html"))
	return &FrontendHandler{
		templates:      tpl,
		postService:    postService,
		sessionService: sessionService,
		commentService: commentService,
	}
}

func (h *FrontendHandler) RegisterFrontendEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.ShowIndex)
	mux.HandleFunc("GET /create-post", h.ShowCreatePost)
	mux.HandleFunc("POST /post", h.CreateNewPost)
	mux.HandleFunc("GET /post/{id}", h.ShowPost)
	// mux.HandleFunc("POST /post", h.CreatePost)
	// mux.HandleFunc("GET /upload", h.ShowUploadForm)
	mux.HandleFunc("GET /archive", h.ShowArchive)
}

func (h *FrontendHandler) ShowIndex(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetActivePosts(r.Context())
	if err != nil {
		h.renderErrorPage(w, err)
	}

	// Get session from context
	session, _ := r.Context().Value("session").(*domain.Session)

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

	// fmt.Println(session)

	h.renderTemplate(w, "archive.html", struct {
		Session *domain.Session
		Posts   []domain.Post
	}{
		Session: session,
		Posts:   posts,
	})
}

func (h *FrontendHandler) ShowCreatePost(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("session").(*domain.Session)
	_ = session

	h.renderTemplate(w, "create-post.html", struct {
	}{})
}

func (h *FrontendHandler) CreateNewPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.renderErrorPage(w, err)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	file, _, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		h.renderErrorPage(w, err)
		return
	}
	if file != nil {
		defer file.Close()
	}

	var imageData []byte
	if file != nil {
		imageData, err = io.ReadAll(file)
		if err != nil {
			h.renderErrorPage(w, err)
			return
		}
	}

	session, ok := r.Context().Value("session").(*domain.Session)
	if !ok {
		h.renderErrorPage(w, err)
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
		h.renderErrorPage(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

func (h *FrontendHandler) ShowPost(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.renderErrorPage(w, err)
		return
	}

	post, err := h.postService.GetPostByID(r.Context(), postID)
	if err != nil {
		h.renderErrorPage(w, err)
		return
	}

	comments, err := CommentService.GetByPostID(h.commentService, r.Context(), postID)
	fmt.Println(comments, err)
	if err != nil {
		h.renderErrorPage(w, err)
	}
	var template string
	if post.IsArchived {
		template = "archive-post.html"
	} else {
		template = "post.html"
	}
	h.renderTemplate(w, template, struct {
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
		httperror.WriteError(w, err)
	}
}

func (h *FrontendHandler) renderErrorPage(w http.ResponseWriter, err error) {
	apiErr := httperror.FromError(err)
	h.renderTemplate(w, "error.html", struct {
		Message    string
		StatusCode int
	}{
		Message:    apiErr.Message,
		StatusCode: apiErr.StatusCode,
	})
}
