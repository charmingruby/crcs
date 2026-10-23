package usecase_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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

type testService struct {
	uc     *usecase.Usecase
	repo   *repository.UserRepository
	tokens *token.Manager
}

func newTestService(t *testing.T) testService {
	t.Helper()

	db := container.StartPostgres(t)

	repo, err := repository.NewUserRepository(db)
	require.NoError(t, err)

	tokens := token.New(testJWTSecret, time.Hour)

	return testService{
		uc:     usecase.New(repo, tokens),
		repo:   repo,
		tokens: tokens,
	}
}

func signUp(t *testing.T, uc *usecase.Usecase) usecase.SignUpOutput {
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
