package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ExonegeS/go-hex-forum/internal/application"
	"github.com/ExonegeS/go-hex-forum/internal/utils"
)

type XHandler struct {
	service *application.XService
	logger  *slog.Logger
}

func NewXHandler(service *application.XService, logger *slog.Logger) *XHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &XHandler{service: service, logger: logger}
}

func (m *XHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/x", m.createX)
}

type createXRequest struct {
	Data string `json:"data"`
}

type createXResponse struct {
	ID int `json:"id"`
}

func (r *createXRequest) Validate() error {
	if r.Data == "" {
		return fmt.Errorf("missing title")
	}
	return nil
}

func (h *XHandler) createX(w http.ResponseWriter, r *http.Request) {
	var req createXRequest
	json.NewDecoder(r.Body).Decode(&req)
	if err := req.Validate(); err != nil {
		h.logger.Error("Failed to validate request", slog.String("error", err.Error()))
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	id, err := h.service.CreateX(r.Context(), req.Data)
	if err != nil {
		h.logger.Error("Failed to create X", slog.String("error", err.Error()))
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp := createXResponse{ID: id}
	utils.WriteJSON(w, http.StatusCreated, resp)
}
