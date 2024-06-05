package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/frangklynndruru/premily_backend/app/controllers/auth"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte(os.Getenv("SECRET_KEY"))

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

		token, err := jwt.ParseWithClaims(tokenCode, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtSecretKey, nil
        })
        if err != nil || !token.Valid {
            // Token tidak valid atau sudah kedaluwarsa
            http.Error(w, "Token has expired or invalid", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        handler.ServeHTTP(w, r.WithContext(ctx))

	})
}
