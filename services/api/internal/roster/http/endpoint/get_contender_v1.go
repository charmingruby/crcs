package endpoint

import (
	"crcs/platform/httpx"
	"net/http"
	"time"

	"crcs/core/model"

	"github.com/charmingruby/crcs/services/api/internal/roster/usecase"
)

type CartelV1 struct {
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            *time.Time `json:"updated_at"`
	LastFightDate        *time.Time `json:"last_fight_date"`
	AvgPerformanceRating *int       `json:"avg_performance_rating"`
	ContenderID          string     `json:"contender_id"`
	WinsByDec            int        `json:"wins_by_dec"`
	WinsByTKO            int        `json:"wins_by_tko"`
	WinsBySub            int        `json:"wins_by_sub"`
	WinsByKO             int        `json:"wins_by_ko"`
	LossesByKO           int        `json:"losses_by_ko"`
	LossesByTKO          int        `json:"losses_by_tko"`
	LossesBySub          int        `json:"losses_by_sub"`
	LossesByDec          int        `json:"losses_by_dec"`
	NoContests           int        `json:"no_contests"`
	Draws                int        `json:"draws"`
	TotalFights          int        `json:"total_fights"`
	Losses               int        `json:"losses"`
	Wins                 int        `json:"wins"`
}

type GetContenderV1Response struct {
	Contender ContenderV1 `json:"contender"`
	Cartel    *CartelV1   `json:"cartel"`
}

func (e *Endpoint) GetContenderV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := httpx.GetPathParam(r, "id")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	contenderID, err := httpx.GetPathParam(r, "cid")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	output, err := e.uc.GetContender(ctx, usecase.GetContenderInput{
		SeasonID:    seasonID,
		ContenderID: contenderID,
	})
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteOKResponse(w, GetContenderV1Response{
		Contender: contenderFromModel(output.Contender),
		Cartel:    cartelFromModel(output.Cartel),
	})
}

func cartelFromModel(cartel *model.Cartel) *CartelV1 {
	if cartel == nil {
		return nil
	}

	return &CartelV1{
		ContenderID:          cartel.ContenderID,
		Wins:                 cartel.Wins,
		Losses:               cartel.Losses,
		Draws:                cartel.Draws,
		NoContests:           cartel.NoContests,
		WinsByKO:             cartel.WinsByKO,
		WinsByTKO:            cartel.WinsByTKO,
		WinsBySub:            cartel.WinsBySub,
		WinsByDec:            cartel.WinsByDec,
		LossesByKO:           cartel.LossesByKO,
		LossesByTKO:          cartel.LossesByTKO,
		LossesBySub:          cartel.LossesBySub,
		LossesByDec:          cartel.LossesByDec,
		AvgPerformanceRating: cartel.AvgPerformanceRating,
		LastFightDate:        cartel.LastFightDate,
		TotalFights:          cartel.TotalFights,
		CreatedAt:            cartel.CreatedAt,
		UpdatedAt:            cartel.UpdatedAt,
	}
}
