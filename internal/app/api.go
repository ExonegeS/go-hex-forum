package apiserver

import (
	"context"
	"database/sql"
	"fmt"
	"go-hex-forum/config"
	"go-hex-forum/internal/adapters/postgres"
	rickmorty "go-hex-forum/internal/adapters/rickmorty_client"
	"go-hex-forum/internal/adapters/storage"
	"go-hex-forum/internal/core/service"
	"go-hex-forum/internal/ports/http/handlers"
	"go-hex-forum/internal/ports/http/middleware"
	"go-hex-forum/internal/utils"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"
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
		"formatTime": utils.FormatTime,
	})

	tpl, err := tpl.ParseGlob(filepath.Join("web", "templates", "*.html"))
	if err != nil {
		s.logger.Error("Failed to load templates", "error", err.Error())
		return err
	}

	//separating frontend handlers from api handlers
	frontendMux := http.NewServeMux()
	apiMux := http.NewServeMux()

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
	frontendHandler := handlers.NewFrontendHandler(PostService, SessionService, CommentService, tpl)
	frontendHandler.RegisterFrontendEndpoints(frontendMux)

	// Middlewares
	SessionMiddleware := SessionHandler.WithSessionToken(int64(s.cfg.SessionConfig.DefaultTTL.Seconds()))
	TimeoutMW := middleware.NewTimeoutContextMW(15)
	MWChain := middleware.NewMiddlewareChain(TimeoutMW, SessionMiddleware, SessionHandler.RequireValidSession)

	serverAddress := fmt.Sprintf("%s:%s", s.cfg.Server.Address, s.cfg.Server.Port)
	s.logger.Info("starting server", slog.String("host", serverAddress))
	httpServer := http.Server{
		Addr:    serverAddress,
		Handler: MWChain(router),
	}
	return httpServer.ListenAndServe()
}
