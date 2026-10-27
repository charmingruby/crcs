package endpoint

import (
	"crcs/platform/httpx"
	"net/http"

	"crcs/core/model/enum"

	"github.com/charmingruby/crcs/services/api/internal/matchmaking/usecase"
)

type RecordMatchResultV1Request struct {
	WinnerID        string `json:"winner_id"        validate:"omitempty,min=1"`
	Outcome         string `json:"outcome"          validate:"required,oneof=W D NC"`
	Method          string `json:"method"           validate:"required,oneof=SUB KO TKO DEC NOC"`
	FightDuration   int    `json:"fight_duration"   validate:"required,gt=0"`
	ScheduledRounds int    `json:"scheduled_rounds" validate:"required,gt=0"`
	ActualRounds    int    `json:"actual_rounds"    validate:"required,gt=0"`
}

type RecordMatchResultV1Response struct {
	ID string `json:"id"`
}

func (e *Endpoint) RecordMatchResultV1(w http.ResponseWriter, r *http.Request) {
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

	request, err := httpx.ParseRequest[RecordMatchResultV1Request](w, r)
	if err != nil {
		return
	}

	output, err := e.uc.RecordMatchResult(ctx, usecase.RecordMatchResultInput{
		SeasonID:        seasonID,
		MatchID:         matchID,
		WinnerID:        request.WinnerID,
		Outcome:         enum.FightOutcome(request.Outcome),
		Method:          enum.FightMethod(request.Method),
		FightDuration:   request.FightDuration,
		ScheduledRounds: request.ScheduledRounds,
		ActualRounds:    request.ActualRounds,
	})
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteCreatedResponse(w, RecordMatchResultV1Response{
		ID: output.ID,
	})
}
