package main

import (
	"database/sql"
	"fmt"
	"go-hex-forum/config"
	apiserver "go-hex-forum/internal/app"
	"go-hex-forum/pkg/lib/prettyslog"
	"html/template"
	"os"
	"path/filepath"
	"time"

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

	// Loading templates for frontend part of the application
	tpl := template.New("post.html").Funcs(template.FuncMap{
		"formatTime": formatTime,
	})

	tpl, err = tpl.ParseGlob(filepath.Join("web", "templates", "*.html"))
	if err != nil {
		logger.Error("Failed to load templates", "error", err.Error())
		return
	}

	// Defining new REST API server
	server := apiserver.NewAPIServer(cfg, db, logger, tpl)
	server.Run()
}

func initDB(cfg *config.Config) (*sql.DB, error) {
	fmt.Println(cfg.DataBase.DBPassword)
	db, err := sql.Open("postgres", cfg.DataBase.MakeConnectionString())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func formatTime(t time.Time) string {
	return t.Format("02 Jan 2006 15:04")
}
