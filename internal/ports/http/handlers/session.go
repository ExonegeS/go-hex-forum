package handlers

import (
	"context"
	"errors"
	"fmt"
	"go-hex-forum/internal/core/domain"
	"go-hex-forum/internal/core/service"
	"go-hex-forum/internal/utils"
	"log"
	"net/http"
	"time"
)

type SessionService interface {
	StoreNewSession() (sessionToken string, err error)
	GetSessionByToken(string) (*domain.Session, error)
}

type SessionHandler struct {
	SessionService SessionService
}

func NewSessionHandler(sessionService SessionService) SessionHandler {
	return SessionHandler{sessionService}
}

func (s *SessionHandler) RegisterEndpoints(mux *http.ServeMux) {
	// mux.HandleFunc("/", helloworld)
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
				token, err := s.SessionService.StoreNewSession()
				fmt.Println(token, "token here")
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
			utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized: no session token"))
			return
		}
		session, err := s.SessionService.GetSessionByToken(cookie.Value)
		fmt.Printf("%#v \n %w", session, err)
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
				utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized: expired session"))
				return
			}
			utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized: invalid session"))
			log.Print(err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
