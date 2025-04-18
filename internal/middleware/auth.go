package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/fehepe/flight-price-service/internal/config"
	"github.com/fehepe/flight-price-service/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
)

// contextKey is a private type to avoid collisions in context.
type contextKey string

const userContextKey = contextKey("userClaims")

// Auth is middleware that enforces a valid JWT in the Authorization header.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			utils.RespondError(w, http.StatusUnauthorized, "missing Authorization header")
			return
		}
		parts := strings.Fields(header)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.RespondError(w, http.StatusUnauthorized, "invalid Authorization header format")
			return
		}

		tokenString := parts[1]
		secret := config.Get("JWT_SECRET", "")
		if secret == "" {
			utils.RespondError(w, http.StatusInternalServerError, "JWT secret not configured")
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), userContextKey, claims))
		next.ServeHTTP(w, r)
	})
}

// FromContext retrieves JWT claims stored in the context.
func FromContext(ctx context.Context) (*jwt.RegisteredClaims, bool) {
	claims, ok := ctx.Value(userContextKey).(*jwt.RegisteredClaims)
	return claims, ok
}
