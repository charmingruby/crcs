package endpoint

import (
	"crcs/platform/httpx"
	"net/http"
)

func (e *Endpoint) GetSeasonV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := httpx.GetPathParam(r, "id")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	season, err := e.uc.GetSeason(ctx, seasonID)
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteOKResponse(w, seasonFromModel(season))
}
