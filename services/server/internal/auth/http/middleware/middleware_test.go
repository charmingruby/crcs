package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/charmingruby/crcs/services/server/internal/auth/http/middleware"
	"github.com/charmingruby/crcs/services/server/internal/auth/token"
)

func TestMiddleware(t *testing.T) {
	tokens := token.New("secret", time.Hour)

	validToken, err := tokens.Sign("user-1")
	require.NoError(t, err)

	expiredToken, err := token.New("secret", -time.Minute).Sign("user-1")
	require.NoError(t, err)

	otherSecretToken, err := token.New("other-secret", time.Hour).Sign("user-1")
	require.NoError(t, err)

	var gotUserID string

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUserID, _ = middleware.UserIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name       string
		header     string
		wantUserID string
		wantStatus int
	}{
		{
			name:       "missing authorization header returns 401",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "malformed authorization header returns 401",
			header:     "Token " + validToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "token signed with another secret returns 401",
			header:     "Bearer " + otherSecretToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "expired token returns 401",
			header:     "Bearer " + expiredToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "valid token passes and exposes the user id",
			header:     "Bearer " + validToken,
			wantStatus: http.StatusOK,
			wantUserID: "user-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID = ""

			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			rec := httptest.NewRecorder()
			middleware.New(tokens)(next).ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantUserID, gotUserID)
		})
	}
}

func TestUserIDFromContextWithoutValue(t *testing.T) {
	userID, ok := middleware.UserIDFromContext(context.Background())

	assert.False(t, ok)
	assert.Empty(t, userID)
}
