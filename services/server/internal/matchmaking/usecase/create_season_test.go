package usecase_test

import (
	"context"
	"crcs/core/customerr"
	"crcs/core/model/enum"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/charmingruby/crcs/services/server/internal/matchmaking/usecase"
)

func TestCreateSeason(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		setup   func(t *testing.T, s testService)
		input   usecase.CreateSeasonInput
		name    string
		errType customerr.ErrorType
		wantErr bool
	}{
		{
			name: "success creates the season as draft",
			input: usecase.CreateSeasonInput{
				Name:    "2026 Spring Season",
				Weeks:   12,
				StartAt: testSeasonStartAt,
			},
		},
		{
			name: "past start_at returns validation error",
			input: usecase.CreateSeasonInput{
				Name:    "2026 Spring Season",
				Weeks:   12,
				StartAt: time.Now().AddDate(-1, 0, 0),
			},
			wantErr: true,
			errType: customerr.TypeValidation,
		},
		{
			name: "duplicate name returns conflict",
			setup: func(t *testing.T, s testService) {
				seedSeason(t, s.uc, "2026 Spring Season")
			},
			input: usecase.CreateSeasonInput{
				Name:    "2026 Spring Season",
				Weeks:   8,
				StartAt: testSeasonStartAt,
			},
			wantErr: true,
			errType: customerr.TypeConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestService(t)

			if tt.setup != nil {
				tt.setup(t, s)
			}

			got, err := s.uc.CreateSeason(ctx, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, customerr.IsType(err, tt.errType))
				assert.Empty(t, got.ID)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, got.ID)

			season, err := s.repo.FindByName(ctx, tt.input.Name)
			require.NoError(t, err)
			require.NotNil(t, season)
			assert.Equal(t, got.ID, season.ID)
			assert.Equal(t, tt.input.Name, season.Name)
			assert.Equal(t, tt.input.Weeks, season.Weeks)
			assert.True(t, season.StartAt.Equal(tt.input.StartAt))
			assert.Equal(t, enum.DraftSeasonStatus, season.Status)
			assert.False(t, season.CreatedAt.IsZero())
		})
	}
}
