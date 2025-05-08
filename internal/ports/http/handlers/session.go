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
	StoreNewSession(context.Context) (sessionToken string, err error)
	GetSessionByToken(context.Context, string) (*domain.Session, error)
	UpdateUserName(ctx context.Context, token string, username string) error
}

type SessionHandler struct {
	SessionService SessionService
}

func NewSessionHandler(sessionService SessionService) SessionHandler {
	return SessionHandler{sessionService}
}

func (s *SessionHandler) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /username", s.UpdateUserName)
}

func (s *SessionHandler) UpdateUserName(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.New("could not parse form"))
	}

	newnickname := r.FormValue("nickname")

	tokenCookie, ok := r.Context().Value("session_token").(string)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	err = s.SessionService.UpdateUserName(r.Context(), tokenCookie, newnickname)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed username change", err))
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
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
			log.Print(err.Error())
			return
		}
		// fmt.Println("Session  here:", session)
		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
