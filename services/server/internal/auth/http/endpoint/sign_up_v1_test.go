package endpoint_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignUpV1(t *testing.T) {
	tests := []struct {
		setup         func(t *testing.T, s testEndpoint)
		body          any
		wantBodyCheck func(t *testing.T, body map[string]any)
		name          string
		wantStatus    int
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
			name: "invalid email returns 400",
			body: map[string]string{
				"first_name": "Jane",
				"last_name":  "Doe",
				"email":      "not-an-email",
				"password":   "password123",
			},
			wantStatus: http.StatusBadRequest,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Contains(t, body["message"], "invalid payload")
			},
		},
		{
			name: "short password returns 400",
			body: map[string]string{
				"first_name": "Jane",
				"last_name":  "Doe",
				"email":      "jane@example.com",
				"password":   "short",
			},
			wantStatus: http.StatusBadRequest,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Contains(t, body["message"], "invalid payload")
			},
		},
		{
			name: "duplicate email returns 409",
			setup: func(t *testing.T, s testEndpoint) {
				seedUser(t, s.uc)
			},
			body: map[string]string{
				"first_name": "Other",
				"last_name":  "User",
				"email":      "jane@example.com",
				"password":   "password123",
			},
			wantStatus: http.StatusConflict,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "email already in use", body["message"])
			},
		},
		{
			name: "success returns 201 and persists the user",
			body: map[string]string{
				"first_name": "Jane",
				"last_name":  "Doe",
				"email":      "jane@example.com",
				"password":   "password123",
			},
			wantStatus: http.StatusCreated,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.NotEmpty(t, body["access_token"])
			},
			wantPersisted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestEndpoint(t)

			if tt.setup != nil {
				tt.setup(t, s)
			}

			rec := doRequest(t, s.ep.SignUpV1, tt.body)

			assert.Equal(t, tt.wantStatus, rec.Code)

			body := decodeBody(t, rec)
			tt.wantBodyCheck(t, body)

			if !tt.wantPersisted {
				return
			}

			accessToken, ok := body["access_token"].(string)
			require.True(t, ok, "expected string access_token in response body")

			userID, err := s.tokens.Verify(accessToken)
			require.NoError(t, err)

			user, err := s.repo.FindByEmail(context.Background(), "jane@example.com")
			require.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, userID, user.ID)
			assert.Equal(t, "Jane", user.FirstName)
		})
	}
}
