package middleware

import (
	"net/http"
)

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool)
	for _, role := range roles {
		allowed[role] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value("role").(string)
			if !ok {
				http.Error(w, "forbidden", 403)
				return
			}
			if !allowed[role] {
				http.Error(w, "permission denied", 403)
				return
			}

			next.ServeHTTP(w, r)
		},
		)
	}
}
