package endpoint_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"crcs/core/model/enum"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSeasonV1(t *testing.T) {
	futureStartAt := func() string {
		return time.Now().AddDate(1, 0, 0).UTC().Format(time.RFC3339)
	}

	validBody := func() map[string]any {
		return map[string]any{
			"name":     "2026 Spring Season",
			"weeks":    12,
			"start_at": futureStartAt(),
		}
	}

	tests := []struct {
		setup         func(t *testing.T, s testEndpoint)
		headers       func(t *testing.T, s testEndpoint) map[string]string
		body          any
		wantBodyCheck func(t *testing.T, body map[string]any)
		name          string
		wantStatus    int
		wantPersisted bool
	}{
		{
			name:       "missing authorization header returns 401",
			body:       validBody(),
			wantStatus: http.StatusUnauthorized,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "missing authorization header", body["message"])
			},
		},
		{
			name: "invalid access token returns 401",
			headers: func(t *testing.T, s testEndpoint) map[string]string {
				return map[string]string{"Authorization": "Bearer not-a-token"}
			},
			body:       validBody(),
			wantStatus: http.StatusUnauthorized,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "invalid access token", body["message"])
			},
		},
		{
			name: "invalid JSON body returns 400",
			headers: func(t *testing.T, s testEndpoint) map[string]string {
				return bearerHeader(t, s.tokenMgr, "user-1")
			},
			body:       "not json",
			wantStatus: http.StatusBadRequest,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Contains(t, body["message"], "invalid payload")
			},
		},
		{
			name: "missing required fields returns 400",
			headers: func(t *testing.T, s testEndpoint) map[string]string {
				return bearerHeader(t, s.tokenMgr, "user-1")
			},
			body: map[string]any{
				"weeks": 12,
			},
			wantStatus: http.StatusBadRequest,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Contains(t, body["message"], "invalid payload")
			},
		},
		{
			name: "zero weeks returns 400",
			headers: func(t *testing.T, s testEndpoint) map[string]string {
				return bearerHeader(t, s.tokenMgr, "user-1")
			},
			body: map[string]any{
				"name":     "2026 Spring Season",
				"weeks":    0,
				"start_at": "2026-01-04T00:00:00Z",
			},
			wantStatus: http.StatusBadRequest,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Contains(t, body["message"], "invalid payload")
			},
		},
		{
			name: "past start_at returns 422",
			headers: func(t *testing.T, s testEndpoint) map[string]string {
				return bearerHeader(t, s.tokenMgr, "user-1")
			},
			body: map[string]any{
				"name":     "2026 Spring Season",
				"weeks":    12,
				"start_at": time.Now().AddDate(-1, 0, 0).UTC().Format(time.RFC3339),
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "season start_at must be a future date", body["message"])
			},
		},
		{
			name: "duplicate name returns 409",
			setup: func(t *testing.T, s testEndpoint) {
				seedSeason(t, s, "2026 Spring Season")
			},
			headers: func(t *testing.T, s testEndpoint) map[string]string {
				return bearerHeader(t, s.tokenMgr, "user-1")
			},
			body:       validBody(),
			wantStatus: http.StatusConflict,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "season already exists", body["message"])
			},
		},
		{
			name: "success returns 201 and persists the season",
			headers: func(t *testing.T, s testEndpoint) map[string]string {
				return bearerHeader(t, s.tokenMgr, "user-1")
			},
			body:       validBody(),
			wantStatus: http.StatusCreated,
			wantBodyCheck: func(t *testing.T, body map[string]any) {
				assert.NotEmpty(t, body["id"])
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

			headers := map[string]string{}
			if tt.headers != nil {
				headers = tt.headers(t, s)
			}

			rec := doRequest(t, s, http.MethodPost, testSeasonRoute, tt.body, headers)

			assert.Equal(t, tt.wantStatus, rec.Code)

			body := decodeBody(t, rec)
			tt.wantBodyCheck(t, body)

			if !tt.wantPersisted {
				return
			}

			id, ok := body["id"].(string)
			require.True(t, ok, "expected string id in response body")

			season, err := s.repo.FindByName(context.Background(), "2026 Spring Season")
			require.NoError(t, err)
			require.NotNil(t, season)
			assert.Equal(t, id, season.ID)
			assert.Equal(t, enum.DraftSeasonStatus, season.Status)
		})
	}
}
