package endpoint

import (
	"crcs/platform/httpx"
	"net/http"
	"time"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/usecase"
)

type CreateSeasonV1Request struct {
	StartAt time.Time `json:"start_at" validate:"required"`
	Name    string    `json:"name"     validate:"required,min=1"`
	Weeks   int       `json:"weeks"    validate:"required,min=1"`
}

type CreateSeasonV1Response struct {
	ID string `json:"id"`
}

func (e *Endpoint) CreateSeasonV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := httpx.ParseRequest[CreateSeasonV1Request](w, r)
	if err != nil {
		return
	}

	output, err := e.uc.CreateSeason(ctx, usecase.CreateSeasonInput{
		Name:    request.Name,
		Weeks:   request.Weeks,
		StartAt: request.StartAt,
	})
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteCreatedResponse(w, CreateSeasonV1Response{
		ID: output.ID,
	})
}
