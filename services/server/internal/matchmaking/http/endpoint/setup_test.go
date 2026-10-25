package endpoint_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/charmingruby/crcs/services/server/internal/matchmaking/http/endpoint"
	"github.com/charmingruby/crcs/services/server/internal/matchmaking/repository"
	"github.com/charmingruby/crcs/services/server/internal/matchmaking/usecase"
	"github.com/charmingruby/crcs/services/server/internal/platform/auth"
	"github.com/charmingruby/crcs/services/server/internal/platform/logging"
	"github.com/charmingruby/crcs/services/server/test/container"
)

const (
	testJWTSecret   = "test-secret"
	testSeasonRoute = "/v1/seasons"
)

var testSeasonStartAt = time.Now().AddDate(1, 0, 0)

func TestMain(m *testing.M) {
	logging.InitLogger()
	os.Exit(m.Run())
}

type testEndpoint struct {
	router   *chi.Mux
	uc       *usecase.Usecase
	repo     *repository.SeasonRepository
	tokenMgr *auth.TokenManager
}

func newTestEndpoint(t *testing.T) testEndpoint {
	t.Helper()

	db := container.StartPostgres(t)

	repo, err := repository.NewSeasonRepository(db)
	require.NoError(t, err)

	tokenMgr := auth.NewTokenManager(testJWTSecret, time.Hour)
	uc := usecase.New(repo)

	router := chi.NewRouter()
	endpoint.New(router, uc, tokenMgr)

	return testEndpoint{
		router:   router,
		uc:       uc,
		repo:     repo,
		tokenMgr: tokenMgr,
	}
}

func seedSeason(t *testing.T, s testEndpoint, name string) usecase.CreateSeasonOutput {
	t.Helper()

	output, err := s.uc.CreateSeason(context.Background(), usecase.CreateSeasonInput{
		Name:    name,
		Weeks:   12,
		StartAt: testSeasonStartAt,
	})
	require.NoError(t, err)

	return output
}

func bearerHeader(t *testing.T, tokenMgr *auth.TokenManager, userID string) map[string]string {
	t.Helper()

	token, err := tokenMgr.Sign(userID)
	require.NoError(t, err)

	return map[string]string{"Authorization": "Bearer " + token}
}

func doRequest(
	t *testing.T,
	s testEndpoint,
	method,
	path string,
	body any,
	headers map[string]string,
) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	if body != nil {
		var bodyBytes []byte

		if raw, ok := body.(string); ok {
			bodyBytes = []byte(raw)
		} else {
			encoded, err := json.Marshal(body)
			require.NoError(t, err)

			bodyBytes = encoded
		}

		reader = bytes.NewReader(bodyBytes)
	}

	req := httptest.NewRequestWithContext(context.Background(), method, path, reader)
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	return rec
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	return body
}
