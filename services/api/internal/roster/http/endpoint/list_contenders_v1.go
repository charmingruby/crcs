package endpoint

import (
	"crcs/platform/httpx"
	"net/http"
	"time"

	"crcs/core/model"

	"github.com/charmingruby/crcs/services/api/internal/roster/usecase"
)

type ContenderV1 struct {
	CreatedAt          time.Time  `json:"created_at"`
	Birthdate          time.Time  `json:"birthdate"`
	Gender             *string    `json:"gender"`
	UpdatedAt          *time.Time `json:"updated_at"`
	Nickname           *string    `json:"nickname"`
	CurrentWeightClass string     `json:"current_weight_class"`
	ID                 string     `json:"id"`
	PlaceOfBirth       string     `json:"place_of_birth"`
	SeasonID           string     `json:"season_id"`
	Status             string     `json:"status"`
	LastName           string     `json:"last_name"`
	FirstName          string     `json:"first_name"`
	Reach              int        `json:"reach"`
	Height             int        `json:"height"`
	Weight             int        `json:"weight"`
}

func (e *Endpoint) ListContendersV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := httpx.GetPathParam(r, "id")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	limit, offset, err := httpx.Pagination(r)
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	query := r.URL.Query()

	output, err := e.uc.ListContenders(ctx, usecase.ListContendersInput{
		SeasonID:    seasonID,
		Status:      query.Get("status"),
		WeightClass: query.Get("weight_class"),
		Gender:      query.Get("gender"),
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	items := make([]ContenderV1, 0, len(output.Items))
	for _, contender := range output.Items {
		items = append(items, contenderFromModel(contender))
	}

	httpx.WriteOKResponse(w, httpx.Page[ContenderV1]{
		Items:  items,
		Total:  output.Total,
		Limit:  limit,
		Offset: offset,
	})
}

func contenderFromModel(contender *model.Contender) ContenderV1 {
	var gender *string
	if contender.Gender != nil {
		value := string(*contender.Gender)
		gender = &value
	}

	return ContenderV1{
		ID:                 contender.ID,
		FirstName:          contender.FirstName,
		LastName:           contender.LastName,
		Nickname:           contender.Nickname,
		Birthdate:          contender.Birthdate,
		PlaceOfBirth:       contender.PlaceOfBirth,
		CurrentWeightClass: string(contender.CurrentWeightClass),
		Reach:              contender.Reach,
		Height:             contender.Height,
		Weight:             contender.Weight,
		Gender:             gender,
		SeasonID:           contender.SeasonID,
		Status:             string(contender.Status),
		CreatedAt:          contender.CreatedAt,
		UpdatedAt:          contender.UpdatedAt,
	}
}
