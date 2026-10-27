package endpoint

import (
	"crcs/platform/httpx"
	"net/http"
	"time"

	"github.com/charmingruby/crcs/services/api/internal/roster/usecase"
)

type IngestContendersV1CartelMatch struct {
	EventName       string    `json:"event_name"       validate:"required,min=1"`
	EventDate       time.Time `json:"event_date"       validate:"required"`
	Method          string    `json:"method"           validate:"required,oneof=SUB KO TKO DEC NOC"`
	Outcome         string    `json:"outcome"          validate:"required,oneof=W L D NC"`
	OpponentName    string    `json:"opponent_name"    validate:"required,min=1"`
	WeightClass     string    `json:"weight_class"     validate:"required,oneof=flyweight bantamweight featherweight lightweight welterweight middleweight light_heavyweight heavyweight"`
	FightDuration   string    `json:"fight_duration"   validate:"required,min=1"`
	EventRating     int       `json:"event_rating"     validate:"gte=0,lte=100"`
	ScheduledRounds int       `json:"scheduled_rounds" validate:"required,min=1"`
	ActualRounds    int       `json:"actual_rounds"    validate:"required,min=1"`
}

type IngestContendersV1Contender struct {
	FirstName          string                          `json:"first_name"           validate:"required,min=1"`
	LastName           string                          `json:"last_name"            validate:"required,min=1"`
	Nickname           string                          `json:"nickname"             validate:"required,min=1"`
	Birthdate          time.Time                       `json:"birthdate"            validate:"required"`
	PlaceOfBirth       string                          `json:"place_of_birth"       validate:"required,min=1"`
	CurrentWeightClass string                          `json:"current_weight_class" validate:"required,oneof=flyweight bantamweight featherweight lightweight welterweight middleweight light_heavyweight heavyweight"`
	Gender             string                          `json:"gender"               validate:"required,oneof=M F"`
	CartelMatches      []IngestContendersV1CartelMatch `json:"cartel_matches"       validate:"dive"`
	Reach              int                             `json:"reach"                validate:"gt=0,lte=250"`
	Height             int                             `json:"height"               validate:"gt=0,lte=250"`
	Weight             int                             `json:"weight"               validate:"gt=0,lte=250"`
}

type IngestContendersV1Request struct {
	Contenders []IngestContendersV1Contender `json:"contenders" validate:"required,min=1,dive"`
}

type IngestContendersV1Response struct {
	ContenderIDs []string `json:"contender_ids"`
}

func (e *Endpoint) IngestContendersV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := httpx.GetPathParam(r, "id")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	request, err := httpx.ParseRequest[IngestContendersV1Request](w, r)
	if err != nil {
		return
	}

	contenders := make([]usecase.IngestContendersContender, 0, len(request.Contenders))

	for _, c := range request.Contenders {
		matches := make([]usecase.IngestContendersCartelMatch, 0, len(c.CartelMatches))

		for _, m := range c.CartelMatches {
			matches = append(matches, usecase.IngestContendersCartelMatch{
				EventName:       m.EventName,
				EventDate:       m.EventDate,
				Method:          m.Method,
				Outcome:         m.Outcome,
				OpponentName:    m.OpponentName,
				WeightClass:     m.WeightClass,
				FightDuration:   m.FightDuration,
				EventRating:     m.EventRating,
				ScheduledRounds: m.ScheduledRounds,
				ActualRounds:    m.ActualRounds,
			})
		}

		contenders = append(contenders, usecase.IngestContendersContender{
			FirstName:          c.FirstName,
			LastName:           c.LastName,
			Nickname:           c.Nickname,
			Birthdate:          c.Birthdate,
			PlaceOfBirth:       c.PlaceOfBirth,
			CurrentWeightClass: c.CurrentWeightClass,
			Gender:             c.Gender,
			CartelMatches:      matches,
			Reach:              c.Reach,
			Height:             c.Height,
			Weight:             c.Weight,
		})
	}

	if err := e.uc.IngestContenders(ctx, usecase.IngestContendersInput{
		SeasonID:   seasonID,
		Contenders: contenders,
	}); err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteAcceptedResponse(w)
}
