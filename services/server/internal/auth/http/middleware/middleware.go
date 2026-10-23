package middleware

import (
	"context"
	"crcs/core/customerr"
	"net/http"
	"strings"

	"github.com/charmingruby/crcs/services/server/internal/auth/token"
	"github.com/charmingruby/crcs/services/server/internal/platform/httpx"
)

type contextKey struct{}

func New(tokens *token.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := bearerToken(r)
			if !ok {
				httpx.HandleErr(w, r, customerr.Unauthorized("missing authorization header", nil))
				return
			}

			userID, err := tokens.Verify(raw)
			if err != nil {
				httpx.HandleErr(w, r, customerr.Unauthorized("invalid access token", err))
				return
			}

			ctx := context.WithValue(r.Context(), contextKey{}, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(contextKey{}).(string)

	return userID, ok && userID != ""
}

func bearerToken(r *http.Request) (string, bool) {
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}

	return parts[1], true
}
