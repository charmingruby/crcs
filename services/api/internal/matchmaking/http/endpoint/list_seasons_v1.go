package endpoint

import (
	"crcs/platform/httpx"
	"net/http"
	"time"

	"crcs/core/model"
	"crcs/core/model/enum"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/usecase"
)

type SeasonV1 struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Status    string     `json:"status"`
	Weeks     int        `json:"weeks"`
	StartAt   time.Time  `json:"start_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (e *Endpoint) ListSeasonsV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset, err := httpx.Pagination(r)
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	var status *enum.SeasonStatus
	if raw := r.URL.Query().Get("status"); raw != "" {
		value := enum.SeasonStatus(raw)
		status = &value
	}

	output, err := e.uc.ListSeasons(ctx, usecase.ListSeasonsInput{
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	items := make([]SeasonV1, 0, len(output.Items))
	for _, season := range output.Items {
		items = append(items, seasonFromModel(season))
	}

	httpx.WriteOKResponse(w, httpx.Page[SeasonV1]{
		Items:  items,
		Total:  output.Total,
		Limit:  limit,
		Offset: offset,
	})
}

func seasonFromModel(season *model.Season) SeasonV1 {
	return SeasonV1{
		ID:        season.ID,
		Name:      season.Name,
		Status:    string(season.Status),
		Weeks:     season.Weeks,
		StartAt:   season.StartAt,
		CreatedAt: season.CreatedAt,
		UpdatedAt: season.UpdatedAt,
	}
}
