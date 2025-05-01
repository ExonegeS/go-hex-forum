package handlers

import (
	"net/http"
	"time"
)

type SessionService interface {
	StoreNewSession() (sessionToken string, err error)
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
			// Пытаемся прочитать куку
			cookie, err := r.Cookie("session_token")
			if err != nil || cookie.Value == "" {
				// Нет токена — создаём новую сессию
				token, err := s.service.StoreNewSession()
				if err != nil {
					http.Error(w, "failed to create session", http.StatusInternalServerError)
					return
				}
				// Устанавливаем куку с нужными параметрами
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    token,
					Path:     "/",
					Expires:  time.Now().Add(time.Duration(expirationInSec) * time.Second),
					HttpOnly: true,
					Secure:   false, // или true, если HTTPS
					SameSite: http.SameSiteLaxMode,
				})
				// Обновляем request context, чтобы дальше могли взять токен
				r = r.Clone(r.Context())
				r.AddCookie(&http.Cookie{Name: "session_token", Value: token})
			}

			// Передаём управление следующему обработчику
			next.ServeHTTP(w, r)
		})
	}
}
