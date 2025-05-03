package handlers

import (
	"fmt"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/utils"
	"html/template"
	"net/http"
	"path/filepath"
)

type FrontendHandler struct {
	templates      *template.Template
	postService    PostService
	sessionService SessionService
}

func NewFrontendHandler(postService PostService, sessionService SessionService) (*FrontendHandler, error) {
	tpl, err := template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		return nil, err
	}

	// Parse component templates
	tpl.ParseGlob(filepath.Join("templates", "components", "*.html"))
	tpl.ParseGlob(filepath.Join("templates", "posts", "*.html"))

	return &FrontendHandler{
		templates:      tpl,
		postService:    postService,
		sessionService: sessionService,
	}, nil
}

func (h *FrontendHandler) RegisterFrontendEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.ShowIndex)
	// mux.HandleFunc("GET /post/{id}", h.ShowPost)
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

	posts := []*domain.Post{
		&domain.Post{
			ID: 1,
			PostAuthor: domain.UserData{
				ID:        100,
				Name:      "Quantum Rick",
				AvatarURL: "https://rickandmortyapi.com/api/character/avatar/274.jpeg",
			},
			Title:      "Post title 1",
			Content:    "Post 1 Contents",
			IsArchived: false,
		},
		&domain.Post{
			ID: 2,
			PostAuthor: domain.UserData{
				ID:        100,
				Name:      "Quantum Rick",
				AvatarURL: "https://rickandmortyapi.com/api/character/avatar/274.jpeg",
			},
			Title:      "Post title 2",
			Content:    "Post 2 Contents",
			IsArchived: true,
		},
	}

	// Get session from context
	session, _ := r.Context().Value("session").(*domain.Session)

	fmt.Println(session)

	h.renderTemplate(w, "base.html", struct {
		Title   string
		Session *domain.Session
		Posts   []*domain.Post
	}{
		Title:   "Title",
		Session: session,
		Posts:   posts,
	})
}

func (h *FrontendHandler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}
