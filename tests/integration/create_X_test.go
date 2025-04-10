package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ExonegeS/go-hex-forum/internal/adapter/inbound/http/handler"
	"github.com/ExonegeS/go-hex-forum/internal/adapter/outbound/repository/postgres"
	"github.com/ExonegeS/go-hex-forum/internal/application"
	"github.com/ExonegeS/go-hex-forum/internal/config"
	"github.com/ExonegeS/go-hex-forum/pkg/lib/prettyslog"

	_ "github.com/lib/pq"
)

func TestHTTPCreateX(t *testing.T) {
	cfg := config.Load("../../.env.test")
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.DBHost,
		cfg.DB.DBPort,
		cfg.DB.DBUser,
		cfg.DB.DBPassword,
		cfg.DB.DBName,
	))
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger := prettyslog.SetupPrettySlog(os.Stdout)

	repo := postgres.NewXRepository(db)
	service := application.NewXService(repo)
	XHandler := handler.NewXHandler(service, logger)

	reqBody, _ := json.Marshal(map[string]string{"data": "integration test"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/x", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	mux := http.NewServeMux()
	XHandler.RegisterRoutes(mux)

	mux.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", res.StatusCode)
	}
}
