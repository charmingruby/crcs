package endpoint_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/charmingruby/crcs/services/server/internal/auth/http/endpoint"
	"github.com/charmingruby/crcs/services/server/internal/auth/repository"
	"github.com/charmingruby/crcs/services/server/internal/auth/token"
	"github.com/charmingruby/crcs/services/server/internal/auth/usecase"
	"github.com/charmingruby/crcs/services/server/internal/platform/logging"
	"github.com/charmingruby/crcs/services/server/test/container"
)

const testJWTSecret = "test-secret"

func TestMain(m *testing.M) {
	logging.InitLogger()
	os.Exit(m.Run())
}

type testEndpoint struct {
	ep     *endpoint.Endpoint
	uc     *usecase.Usecase
	repo   *repository.UserRepository
	tokens *token.Manager
}

func newTestEndpoint(t *testing.T) testEndpoint {
	t.Helper()

	db := container.StartPostgres(t)

	repo, err := repository.NewUserRepository(db)
	require.NoError(t, err)

	tokens := token.New(testJWTSecret, time.Hour)
	uc := usecase.New(repo, tokens)

	return testEndpoint{
		ep:     endpoint.New(chi.NewRouter(), uc),
		uc:     uc,
		repo:   repo,
		tokens: tokens,
	}
}

func seedUser(t *testing.T, uc *usecase.Usecase) usecase.SignUpOutput {
	t.Helper()

	output, err := uc.SignUp(context.Background(), usecase.SignUpInput{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane@example.com",
		Password:  "password123",
	})
	require.NoError(t, err)

	return output
}

func doRequest(t *testing.T, handler http.HandlerFunc, body any) *httptest.ResponseRecorder {
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

	req := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/", reader)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler(rec, req)

	return rec
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

	return body
}
