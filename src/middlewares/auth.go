package middlewares

import (
	"context"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		auth := r.Header.Get("Authorization")
		if auth != "" {
			ctx = context.WithValue(ctx, "userID", auth)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
