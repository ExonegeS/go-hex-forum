package handlers

import (
	"context"
	"errors"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/core/service"
	"go-hex-forum/internal/ports/http/httperror"
	"go-hex-forum/internal/utils"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

type SessionService interface {
	StoreNewSession(context.Context) (sessionToken string, err error)
	GetSessionByToken(context.Context, string) (*domain.Session, error)
	UpdateUserName(ctx context.Context, token string, username string) error
}

type SessionHandler struct {
	templates      *template.Template
	SessionService SessionService
	logger         *slog.Logger
}

func NewSessionHandler(tmpl *template.Template, sessionService SessionService, logger *slog.Logger) SessionHandler {
	return SessionHandler{tmpl, sessionService, logger}
}

func (s *SessionHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /username", s.UpdateUserName)
}

func (s *SessionHandler) UpdateUserName(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		httperror.WriteError(w, errors.New("could not parse form"))
	}

	newnickname := r.FormValue("nickname")
	source := r.FormValue("source")
	tokenCookie, ok := r.Context().Value("session_token").(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	err = s.SessionService.UpdateUserName(r.Context(), tokenCookie, newnickname)
	if err != nil {
		if source == "frontend" {
			fmt.Println("here")
			s.renderErrorPage(w, err)
			return
		}
		httperror.WriteError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *SessionHandler) WithSessionToken(expirationInSec int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("session_token")
			var token string
			if err != nil || tokenCookie.Value == "" {
				token, err = s.SessionService.StoreNewSession(r.Context())
				if err != nil {
					utils.WriteError(w, http.StatusInternalServerError, errors.New("failed to create session"))
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    token,
					Expires:  time.Now().Add(time.Duration(expirationInSec) * time.Second),
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteLaxMode,
					Path:     "/",
				})
			} else {
				token = tokenCookie.Value
			}

			ctx := context.WithValue(r.Context(), "session_token", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (s *SessionHandler) RequireValidSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		token := cookie.Value
		if err != nil || token == "" {
			ctxSessionToken, ok := r.Context().Value("session_token").(string)
			if !ok || ctxSessionToken == "" {
				utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized: no session token"))
				return
			}
			token = ctxSessionToken
		}
		session, err := s.SessionService.GetSessionByToken(r.Context(), token)
		if err != nil || session == nil {
			if errors.Is(err, service.ErrSessionExpired) {
				http.SetCookie(w, &http.Cookie{
					Name:   "session_token",
					Value:  "",
					MaxAge: -1,
				})
				utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized: expired session"))
				return
			}
			utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized: invalid session"))
			return
		}
		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *SessionHandler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	err := h.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		httperror.WriteError(w, err)
	}
}

func (h *SessionHandler) renderErrorPage(w http.ResponseWriter, err error) {
	apiErr := httperror.FromError(err)
	h.renderTemplate(w, "error.html", struct {
		Message    string
		StatusCode int
	}{
		Message:    apiErr.Message,
		StatusCode: apiErr.StatusCode,
	})
}
