package middleware

import (
	"context"
	"net/http"
	"time"
)

type Middleware func(next http.Handler) http.Handler

func NewMiddlewareChain(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

func NewTimeoutContextMW(timeoutInSec int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutInSec))
				defer cancel()

				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			})
	}
}

func RecoveryMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
