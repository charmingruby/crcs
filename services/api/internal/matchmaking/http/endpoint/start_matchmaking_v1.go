package endpoint

import (
	"net/http"

	"crcs/platform/httpx"
)

func (e *Endpoint) StartMatchmakingV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := httpx.GetPathParam(r, "id")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	if err := e.uc.StartMatchmaking(ctx, seasonID); err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteAcceptedResponse(w)
}
