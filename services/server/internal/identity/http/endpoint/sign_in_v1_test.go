package endpoint_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignInV1(t *testing.T) {
	tests := []struct {
		body          any
		wantBodyCheck func(t *testing.T, body map[string]any)
		name          string
		wantStatus    int
		wantToken     bool
		wantPersisted bool
	}{
		{
			name:       "invalid JSON body returns 400",
			body:       "not json",
			wantStatus: http.StatusBadRequest,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Contains(t, body["message"], "invalid payload")
			},
		},
		{
			name: "missing required fields returns 400",
			body: map[string]string{
				"email": "jane@example.com",
			},
			wantStatus: http.StatusBadRequest,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Contains(t, body["message"], "invalid payload")
			},
		},
		{
			name: "unknown email returns 401",
			body: map[string]string{
				"email":    "unknown@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusUnauthorized,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "invalid credentials", body["message"])
			},
		},
		{
			name: "wrong password returns 401",
			body: map[string]string{
				"email":    "jane@example.com",
				"password": "wrong-password",
			},
			wantStatus: http.StatusUnauthorized,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "invalid credentials", body["message"])
			},
		},
		{
			name: "success returns 200 and a valid token",
			body: map[string]string{
				"email":    "jane@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusOK,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.NotEmpty(t, body["access_token"])
			},
			wantToken:     true,
			wantPersisted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestEndpoint(t)

			if tt.wantPersisted {
				seedUser(t, s.uc)
			}

			rec := doRequest(t, s.ep.SignInV1, tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)

			body := decodeBody(t, rec)
			tt.wantBodyCheck(t, body)

			if !tt.wantToken {
				return
			}

			accessToken, ok := body["access_token"].(string)
			require.True(t, ok, "expected string access_token in response body")

			userID, err := s.tokenMgr.Verify(accessToken)
			require.NoError(t, err)

			user, err := s.repo.FindByEmail(context.Background(), "jane@example.com")
			require.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, userID, user.ID)
		})
	}
}
