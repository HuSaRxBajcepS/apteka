package middleware

import (
	"apteka/internal/auth"
	"context"
	"net/http"
	"strings"
)

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, "missing token", 401)
				return
			}
			parts := strings.Split(header, " ")
			if len(parts) != 2 {
				http.Error(w, "invalid token", 401)
				return
			}
			claims, err := auth.
				ParseToken(parts[1], secret)
			if err != nil {
				http.Error(w, "unauthorized", 401)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			ctx = context.WithValue(ctx, "role", claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
		)
	}
}
