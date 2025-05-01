package handlers

import (
	"context"
	"errors"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/core/service"
	"log"
	"net/http"
	"time"
)

type SessionService interface {
	StoreNewSession() (sessionToken string, err error)
	GetSessionByToken(string) (*domain.Session, error)
}

type SessionHandler struct {
	service SessionService
}

func NewSessionHandler(sessionService SessionService) SessionHandler {
	return SessionHandler{sessionService}
}

func (s *SessionHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/", helloworld)
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func (s *SessionHandler) WithSessionToken(expirationInSec int64) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			var token string
			if err != nil || cookie.Value == "" {
				token, err := s.service.StoreNewSession()
				if err != nil {
					http.Error(w, "failed to create session", http.StatusInternalServerError)
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
				token = cookie.Value
			}

			ctx := context.WithValue(r.Context(), "session_token", token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (s *SessionHandler) RequireValidSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "unauthorized: no session token", http.StatusUnauthorized)
			return
		}

		session, err := s.service.GetSessionByToken(cookie.Value)
		if err != nil || session == nil {
			if errors.Is(err, service.ErrSessionExpired) {
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    "",
					Path:     "/",
					Expires:  time.Unix(0, 0),
					MaxAge:   -1,
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteLaxMode,
				})
				http.Error(w, "unauthorized: expired session", http.StatusUnauthorized)
				return
			}
			http.Error(w, "unauthorized: invalid session", http.StatusUnauthorized)
			log.Print(err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
