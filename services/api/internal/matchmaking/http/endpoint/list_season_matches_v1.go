package endpoint

import (
	"crcs/platform/httpx"
	"net/http"
	"time"

	"crcs/core/model"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/usecase"
)

type MatchResultV1 struct {
	CreatedAt       time.Time `json:"created_at"`
	WinnerID        *string   `json:"winner_id"`
	ID              string    `json:"id"`
	MatchID         string    `json:"match_id"`
	Outcome         string    `json:"outcome"`
	Method          string    `json:"method"`
	FightDuration   int       `json:"fight_duration"`
	ScheduledRounds int       `json:"scheduled_rounds"`
	ActualRounds    int       `json:"actual_rounds"`
}

type MatchV1 struct {
	CreatedAt           time.Time      `json:"created_at"`
	Gender              *string        `json:"gender"`
	Result              *MatchResultV1 `json:"result"`
	UpdatedAt           *time.Time     `json:"updated_at"`
	ID                  string         `json:"id"`
	SeasonID            string         `json:"season_id"`
	FavoriteContenderID string         `json:"favorite_contender_id"`
	UnderdogContenderID string         `json:"underdog_contender_id"`
	WeightClass         string         `json:"weight_class"`
	Week                int            `json:"week"`
	Position            int            `json:"position"`
}

func (e *Endpoint) ListSeasonMatchesV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := httpx.GetPathParam(r, "id")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	week, err := httpx.QueryInt(r, "week", 0)
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	output, err := e.uc.ListSeasonMatches(ctx, usecase.ListSeasonMatchesInput{
		SeasonID: seasonID,
		Week:     week,
	})
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	items := make([]MatchV1, 0, len(output.Items))
	for _, item := range output.Items {
		items = append(items, matchFromModel(item))
	}

	httpx.WriteOKResponse(w, httpx.Page[MatchV1]{
		Items:  items,
		Total:  len(items),
		Limit:  0,
		Offset: 0,
	})
}

func matchFromModel(item usecase.SeasonMatch) MatchV1 {
	var gender *string
	if item.Match.Gender != nil {
		value := string(*item.Match.Gender)
		gender = &value
	}

	return MatchV1{
		ID:                  item.Match.ID,
		SeasonID:            item.Match.SeasonID,
		Week:                item.Match.Week,
		Position:            item.Match.Position,
		FavoriteContenderID: item.Match.FavoriteContenderID,
		UnderdogContenderID: item.Match.UnderdogContenderID,
		WeightClass:         string(item.Match.WeightClass),
		Gender:              gender,
		Result:              matchResultFromModel(item.Result),
		CreatedAt:           item.Match.CreatedAt,
		UpdatedAt:           item.Match.UpdatedAt,
	}
}

func matchResultFromModel(result *model.MatchResult) *MatchResultV1 {
	if result == nil {
		return nil
	}

	return &MatchResultV1{
		ID:              result.ID,
		MatchID:         result.MatchID,
		WinnerID:        result.WinnerID,
		Outcome:         string(result.Outcome),
		Method:          string(result.Method),
		FightDuration:   result.FightDuration,
		ScheduledRounds: result.ScheduledRounds,
		ActualRounds:    result.ActualRounds,
		CreatedAt:       result.CreatedAt,
	}
}
