package main

import (
	"database/sql"
	"go-hex-forum/config"
	apiserver "go-hex-forum/internal/app"
	"go-hex-forum/pkg/lib/prettyslog"

	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Load the configuration
	cfg := config.NewConfig()

	// Initialize the logger
	logger := prettyslog.SetupPrettySlog(os.Stdout)

	// Initialize the database object and ping the database
	db, err := initDB(cfg)
	if err != nil {
		logger.Error("Failed to connect to the database", "error", err.Error())
		return
	}
	defer db.Close()

	// Defining new REST API server
	server := apiserver.NewAPIServer(cfg, db, logger)
	server.Run()
}

func initDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DataBase.MakeConnectionString())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
