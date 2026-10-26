package public

import (
	"context"
	"crcs/core/customerr"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/usecase"
)

type SeasonReader struct {
	uc *usecase.Usecase
}

func NewSeasonReader(uc *usecase.Usecase) *SeasonReader {
	return &SeasonReader{
		uc: uc,
	}
}

func (s *SeasonReader) SeasonExists(ctx context.Context, id string) (bool, error) {
	season, err := s.uc.GetSeason(ctx, id)

	if err != nil {
		if customerr.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return season != nil, nil
}
