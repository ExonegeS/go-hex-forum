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
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

type APIServer struct {
	cfg    *config.Config
	db     *sql.DB
	logger *slog.Logger
	tpl    *template.Template
}

func NewAPIServer(config *config.Config, db *sql.DB, logger *slog.Logger, tpl *template.Template) *APIServer {
	return &APIServer{config, db, logger, tpl}
}

func (s *APIServer) Run() error {
	ctx := context.Background()

	frontendHandlers := http.NewServeMux()
	apiHandlers := http.NewServeMux()
	router := http.NewServeMux()

	router.Handle("/api/", http.StripPrefix("/api", apiHandlers))
	router.Handle("/", frontendHandlers)

	// Transactor
	transactor := postgres.NewTransactor(s.db)

	// Third Party APIs
	UserdataProvider := rickmorty.NewUserDataProvider("https://rickandmortyapi.com/api", 826)
	ImageStorage := storage.NewImageStorage(s.cfg.Storage.MakeAddressString(), s.cfg.Storage.MaxNameLength)

	// Session
	SessionRepository := postgres.NewSessionRepository(s.db)
	SessionService := service.NewSessionService(SessionRepository, time.Now, UserdataProvider, s.cfg.SessionConfig)
	SessionHandler := handlers.NewSessionHandler(SessionService)
	SessionHandler.RegisterEndpoints(apiHandlers)

	// Post
	PostRepository := postgres.NewPostRepository(s.db)
	PostService := service.NewPostService(PostRepository, ImageStorage)

	go PostService.ArchiveExpiredPostsWorker(ctx)
	PostHandler := handlers.NewPostHandler(PostService)
	PostHandler.RegisterEndpoints(apiHandlers)

	// Comment
	CommentRepository := postgres.NewCommentRepository(s.db)
	CommentService := service.NewCommentService(transactor, CommentRepository, PostRepository, ImageStorage)
	CommentHandler := handlers.NewCommentHandler(CommentService)
	CommentHandler.RegisterEndpoints(apiHandlers)

	// rendering pages and calls on /api
	frontendHandler := handlers.NewFrontendHandler(PostService, SessionService, CommentService, s.tpl)
	frontendHandler.RegisterFrontendEndpoints(frontendHandlers)

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
