package apiserver

import (
	"database/sql"
	"fmt"
	"go-hex-forum/config"
	"go-hex-forum/internal/adapters/postgres"
	"go-hex-forum/internal/core/service"
	"go-hex-forum/internal/ports/http/handlers"
	"go-hex-forum/internal/ports/http/middleware"
	"log/slog"
	"net/http"
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
	mux := http.NewServeMux()

	SessionsRepository := postgres.NewSessionRepository(s.db)
	SessionService := service.NewSessionService(SessionsRepository, time.Now, nil, s.cfg.SessionConfig)
	SessionHandler := handlers.NewSessionHandler(SessionService)
	SessionHandler.RegisterEndpoints(mux)

	SessionMiddleware := SessionHandler.WithSessionToken(int64(s.cfg.SessionConfig.DefaultTTL.Seconds()))
	timeoutMW := middleware.NewTimoutContextMW(15)
	MWChain := middleware.NewMiddlewareChain(middleware.RecoveryMW, timeoutMW, SessionMiddleware)

	serverAddress := fmt.Sprintf("%s:%s", s.cfg.Server.Address, s.cfg.Server.Port)
	s.logger.Info("starting server", slog.String("host", serverAddress))
	httpServer := http.Server{
		Addr:    serverAddress,
		Handler: MWChain(mux),
	}
	return httpServer.ListenAndServe()
}
