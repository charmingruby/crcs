package endpoint

import (
	"crcs/platform/httpx"
	"net/http"
	"time"

	"crcs/core/model/enum"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/usecase"
)

type UpdateSeasonV1Request struct {
	Name    *string    `json:"name"     validate:"omitempty,min=1"`
	Status  *string    `json:"status"   validate:"omitempty,oneof=accepted rejected"`
	StartAt *time.Time `json:"start_at" validate:"omitempty"`
}

func (e *Endpoint) UpdateSeasonV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := httpx.GetPathParam(r, "id")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	request, err := httpx.ParseRequest[UpdateSeasonV1Request](w, r)
	if err != nil {
		return
	}

	input := usecase.UpdateSeasonInput{
		ID:      seasonID,
		Name:    request.Name,
		StartAt: request.StartAt,
	}

	if request.Status != nil {
		status := enum.SeasonStatus(*request.Status)
		input.Status = &status
	}

	if err := e.uc.UpdateSeason(ctx, input); err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteNoContentResponse(w)
}
