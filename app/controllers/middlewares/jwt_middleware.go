package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
)

func JwtMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required!", http.StatusUnauthorized)
			return
		}

		tokenCode := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := auth.TokenVerify(tokenCode)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
