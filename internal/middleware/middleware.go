package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/IFA-01/messenger/internal/auth"
)

func AuthMiddleWare(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			user, err := auth.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
