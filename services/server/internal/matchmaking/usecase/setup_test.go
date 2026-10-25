package usecase_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/charmingruby/crcs/services/server/internal/matchmaking/repository"
	"github.com/charmingruby/crcs/services/server/internal/matchmaking/usecase"
	"github.com/charmingruby/crcs/services/server/internal/platform/logging"
	"github.com/charmingruby/crcs/services/server/test/container"
)

var testSeasonStartAt = time.Now().AddDate(1, 0, 0)

func TestMain(m *testing.M) {
	logging.InitLogger()
	os.Exit(m.Run())
}

type testService struct {
	uc   *usecase.Usecase
	repo *repository.SeasonRepository
}

func newTestService(t *testing.T) testService {
	t.Helper()

	db := container.StartPostgres(t)

	repo, err := repository.NewSeasonRepository(db)
	require.NoError(t, err)

	return testService{
		uc:   usecase.New(repo),
		repo: repo,
	}
}

func seedSeason(t *testing.T, uc *usecase.Usecase, name string) usecase.CreateSeasonOutput {
	t.Helper()

	output, err := uc.CreateSeason(context.Background(), usecase.CreateSeasonInput{
		Name:    name,
		Weeks:   12,
		StartAt: testSeasonStartAt,
	})
	require.NoError(t, err)

	return output
}
