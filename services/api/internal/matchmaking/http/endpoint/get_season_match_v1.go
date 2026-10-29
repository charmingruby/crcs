package endpoint

import (
	"crcs/platform/httpx"
	"net/http"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/usecase"
)

func (e *Endpoint) GetSeasonMatchV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := httpx.GetPathParam(r, "id")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	matchID, err := httpx.GetPathParam(r, "mid")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	match, err := e.uc.GetSeasonMatch(ctx, usecase.GetSeasonMatchInput{
		SeasonID: seasonID,
		MatchID:  matchID,
	})
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteOKResponse(w, matchFromModel(match))
}
