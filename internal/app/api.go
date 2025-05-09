package apiserver

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go-hex-forum/config"
	"go-hex-forum/internal/adapters/postgres"
	rickmorty "go-hex-forum/internal/adapters/rickmorty_client"
	"go-hex-forum/internal/adapters/storage"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/core/service"
	"go-hex-forum/internal/ports/http/handlers"
	"go-hex-forum/internal/ports/http/middleware"
	"go-hex-forum/internal/utils"
)

type APIServer struct {
	cfg    *config.Config
	db     *sql.DB
	logger *slog.Logger
}

func NewAPIServer(config *config.Config, db *sql.DB, logger *slog.Logger) *APIServer {
	return &APIServer{config, db, logger}
}

func (s *APIServer) Run() error {
	ctx := context.Background()
	if s.cfg == nil {
		return fmt.Errorf("cannot start, config is nil")
	}
	if s.db == nil {
		return fmt.Errorf("cannot start, db is nil")
	}
	if s.logger == nil {
		return fmt.Errorf("cannot start, logger is nil")
	}

	// Loading templates for frontend part of the application
	tpl := template.New("post.html").Funcs(template.FuncMap{
		"commentArgs": func(comment *domain.Comment, allComments []*domain.Comment, postID int64, depth int) map[string]interface{} {
			return map[string]interface{}{
				"Comment":     comment,
				"AllComments": allComments,
				"PostID":      postID,
				"Depth":       depth,
			}
		},
		"formatTime": utils.FormatTime,
		"getReplies": func(comments []*domain.Comment, parentID int64) []*domain.Comment {
			var replies []*domain.Comment
			for _, c := range comments {
				if c.ParentCommentID != nil && *c.ParentCommentID == parentID {
					replies = append(replies, c)
				}
			}
			return replies
		},
		"add": func(a, b int) int { return a + b },
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.ReplaceAll(template.HTMLEscapeString(text), "\n", "<br>"))
		},
	})

	tpl, err := tpl.ParseGlob(filepath.Join("web", "templates", "*.html"))
	if err != nil {
		s.logger.Error("Failed to load templates", "error", err.Error())
		return err
	}

	// separating frontend handlers from api handlers
	frontendMux := http.NewServeMux()
	apiMux := http.NewServeMux()

	frontendMux.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir(filepath.Join("web", "static"))),
		),
	)
	router := http.NewServeMux()

	router.Handle("/api/", http.StripPrefix("/api", apiMux))
	router.Handle("/", frontendMux)

	// Transactor
	transactor := postgres.NewTransactor(s.db)

	// Third Party APIs
	UserdataProvider := rickmorty.NewUserDataProvider("https://rickandmortyapi.com/api", 826)
	ImageStorage := storage.NewImageStorage(s.cfg.Storage.MakeAddressString(), s.cfg.Storage.MaxNameLength)

	// Session
	SessionRepository := postgres.NewSessionRepository(s.db)
	SessionService := service.NewSessionService(SessionRepository, time.Now, UserdataProvider, s.cfg.SessionConfig)
	SessionHandler := handlers.NewSessionHandler(tpl, SessionService, s.logger)
	SessionHandler.RegisterEndpoints(apiMux)

	// Post
	PostRepository := postgres.NewPostRepository(s.db)
	PostService := service.NewPostService(PostRepository, ImageStorage)

	go PostService.ArchiveExpiredPostsWorker(ctx)
	PostHandler := handlers.NewPostHandler(PostService)
	PostHandler.RegisterEndpoints(apiMux)

	// Comment
	CommentRepository := postgres.NewCommentRepository(s.db)
	CommentService := service.NewCommentService(transactor, CommentRepository, PostRepository, ImageStorage)
	CommentHandler := handlers.NewCommentHandler(CommentService)
	CommentHandler.RegisterEndpoints(apiMux)

	// rendering pages and calls on /api
	frontendHandler := handlers.NewFrontendHandler(PostService, SessionService, CommentService, tpl, s.logger)
	frontendHandler.RegisterFrontendEndpoints(frontendMux)

	// Middlewares
	SessionMiddleware := SessionHandler.WithSessionToken(int64(s.cfg.SessionConfig.DefaultTTL.Seconds()))
	TimeoutMW := middleware.NewTimeoutContextMW(15)
	MWChain := middleware.NewMiddlewareChain(middleware.RecoveryMW, TimeoutMW, SessionMiddleware, SessionHandler.RequireValidSession)

	serverAddress := fmt.Sprintf("%s:%s", s.cfg.Server.Address, s.cfg.Server.Port)
	s.logger.Info("starting server", slog.String("host", serverAddress))
	httpServer := http.Server{
		Addr:    serverAddress,
		Handler: MWChain(router),
	}
	return httpServer.ListenAndServe()
}
