package handlers

import (
	"errors"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/utils"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type FrontendHandler struct {
	templates      *template.Template
	postService    PostService
	sessionService SessionService
}

func NewFrontendHandler(postService PostService, sessionService SessionService) (*FrontendHandler, error) {
	tpl, err := template.ParseGlob(filepath.Join("web", "templates", "*.html"))
	if err != nil {
		return nil, err
	}

	// Parse component templates
	tpl.ParseGlob(filepath.Join("web", "templates", "components", "*.html"))
	// tpl.ParseGlob(filepath.Join("templates", "posts", "*.html"))

	return &FrontendHandler{
		templates:      tpl,
		postService:    postService,
		sessionService: sessionService,
	}, nil
}

func (h *FrontendHandler) RegisterFrontendEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.ShowIndex)
	mux.HandleFunc("GET /create-post.html", h.CreatePost)
	mux.HandleFunc("GET /post/{id}", h.ShowPost)
	// mux.HandleFunc("POST /post", h.CreatePost)
	// mux.HandleFunc("GET /upload", h.ShowUploadForm)
}

func (h *FrontendHandler) ShowIndex(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// Get active posts
	// posts, err := h.postService.GetActivePosts(ctx)
	// if err != nil {
	// 	http.Error(w, "Error loading posts", http.StatusInternalServerError)
	// 	return
	// }

	posts := []domain.Post{
		domain.Post{
			ID: 1,
			PostAuthor: domain.UserData{
				ID:        100,
				Name:      "Quantum Rick",
				AvatarURL: "https://rickandmortyapi.com/api/character/avatar/274.jpeg",
			},
			ImagePath:  "http://localhost:6969/user-27/uKsfIT",
			Title:      "Post title 1",
			Content:    "Post 1 Contents",
			IsArchived: false,
		},
	}

	postsNew, err := h.postService.GetActivePosts(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
	posts = append(posts, postsNew...)

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

	h.renderTemplate(w, "post.html", struct {
		UserAvatar string
		UserName   string
		DataTime   string
		PostID     int64
		ImagePath  string
		Title      string
		Content    string
		Comments   []domain.Comment
	}{
		UserAvatar: post.PostAuthor.AvatarURL,
		UserName:   post.PostAuthor.Name,
		DataTime:   time.Now().Format("02 Jan 2006 15:04"),
		PostID:     post.ID,
		ImagePath:  post.ImagePath,
		Title:      post.Title,
		Content:    post.Content,
		Comments: []domain.Comment{
			{
				ID:              99,
				ParentCommentID: nil,
				Content:         "I hate n",
				ImagePath:       "https://rickandmortyapi.com/api/character/avatar/300.jpeg",
				Author: domain.UserData{
					ID:        93,
					Name:      "Nurs",
					AvatarURL: "https://rickandmortyapi.com/api/character/avatar/174.jpeg",
				},
				CreatedAt: time.Now(),
			},
			{
				ID:              93,
				ParentCommentID: nil,
				Content:         "I hate n",
				ImagePath:       "",
				Author: domain.UserData{
					ID:        93,
					Name:      "Nurs",
					AvatarURL: "https://rickandmortyapi.com/api/character/avatar/174.jpeg",
				},
				CreatedAt: time.Now(),
			},
			{
				ID:              99,
				ParentCommentID: nil,
				Content:         "I hate n",
				ImagePath:       "https://rickandmortyapi.com/api/character/avatar/600.jpeg",
				Author: domain.UserData{
					ID:        93,
					Name:      "Nurs",
					AvatarURL: "https://rickandmortyapi.com/api/character/avatar/174.jpeg",
				},
				CreatedAt: time.Now(),
			},
		},
	})
}

func (h *FrontendHandler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}
