package middleware

import (
	"context"
	"net/http"
	"strings"

	"ridehailing/backend/internal/pkg/jwtutil"
	"ridehailing/backend/internal/pkg/response"
)

type contextKey string

const userContextKey contextKey = "auth_user"

type CurrentAuthUser struct {
	ID    uint   `json:"id"`
	Role  string `json:"role"`
	Phone string `json:"phone"`
}

func AuthMiddleware(jwtManager *jwtutil.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(w, http.StatusUnauthorized, "missing or invalid Authorization header")
				return
			}

			tokenStr := strings.TrimSpace(parts[1])
			if tokenStr == "" {
				response.Error(w, http.StatusUnauthorized, "missing token")
				return
			}

			claims, err := jwtManager.ParseToken(tokenStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}

			currentUser := &CurrentAuthUser{
				ID:    claims.UserID,
				Role:  claims.Role,
				Phone: claims.Phone,
			}

			ctx := context.WithValue(r.Context(), userContextKey, currentUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CurrentUser(r *http.Request) (*CurrentAuthUser, bool) {
	v := r.Context().Value(userContextKey)
	user, ok := v.(*CurrentAuthUser)
	return user, ok
}
