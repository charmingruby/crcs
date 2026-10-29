package on_season_matchmaking_requested_v1

import (
	"math"
	"sort"
	"time"

	"crcs/core/model"
	"crcs/core/model/enum"
)

const maxFightsPerWeek = 5

type seededContender struct {
	contender *model.Contender
	cartel    *model.Cartel
	score     int
}

type fight struct {
	favorite    *seededContender
	underdog    *seededContender
	weightClass enum.WeightClass
	gender      enum.Gender
	appeal      int
}

func buildSchedule(
	season *model.Season,
	contenders []*model.Contender,
	cartels []*model.Cartel,
) ([]*model.Match, error) {
	cartelsByContenderID := make(map[string]*model.Cartel, len(cartels))
	for _, cartel := range cartels {
		cartelsByContenderID[cartel.ContenderID] = cartel
	}

	divisions := make(map[string][]seededContender)
	for _, contender := range contenders {
		if contender.Gender == nil {
			continue
		}

		cartel := cartelsByContenderID[contender.ID]
		key := string(contender.CurrentWeightClass) + "|" + string(*contender.Gender)

		divisions[key] = append(divisions[key], seededContender{
			contender: contender,
			cartel:    cartel,
			score:     scoreContender(cartel),
		})
	}

	fights := make([]fight, 0)
	for _, division := range divisions {
		sortSeeded(division)

		for i := 0; i+1 < len(division); i += 2 {
			favorite, underdog := division[i], division[i+1]

			fights = append(fights, fight{
				favorite:    &favorite,
				underdog:    &underdog,
				weightClass: favorite.contender.CurrentWeightClass,
				gender:      *favorite.contender.Gender,
				appeal:      fightAppeal(favorite.score, underdog.score),
			})
		}
	}

	sortFights(fights)

	cards := distributeWeeks(fights, season.Weeks)

	matches := make([]*model.Match, 0, len(fights))
	for weekIndex, card := range cards {
		for position, fight := range orderCard(card) {
			match, err := model.NewMatch(model.MatchInput{
				SeasonID:            season.ID,
				Week:                weekIndex + 1,
				Position:            position + 1,
				FavoriteContenderID: fight.favorite.contender.ID,
				UnderdogContenderID: fight.underdog.contender.ID,
				WeightClass:         fight.weightClass,
				Gender:              fight.gender,
			})
			if err != nil {
				return nil, err
			}

			matches = append(matches, match)
		}
	}

	return matches, nil
}

// scoreContender ranks a contender from 0 to 100, weighing record, performance
// quality and experience. A missing cartel is scored as a newcomer.
func scoreContender(cartel *model.Cartel) int {
	record := 0.0
	quality := 50.0
	totalFights := 0

	if cartel != nil {
		totalFights = cartel.TotalFights

		if totalFights > 0 {
			record = (float64(cartel.Wins) + 0.5*float64(cartel.Draws)) / float64(totalFights)
		}

		if cartel.AvgPerformanceRating != nil {
			quality = float64(*cartel.AvgPerformanceRating)
		}
	}

	experience := float64(minInt(totalFights, 10)) * 10
	score := 0.45*record*100 + 0.35*quality + 0.20*experience

	return clampScore(int(math.Round(score)))
}

func sortSeeded(division []seededContender) {
	sort.SliceStable(division, func(i, j int) bool {
		if division[i].score != division[j].score {
			return division[i].score > division[j].score
		}

		iFights, jFights := totalFights(division[i].cartel), totalFights(division[j].cartel)
		if iFights != jFights {
			return iFights > jFights
		}

		iLast, jLast := lastFightDate(division[i].cartel), lastFightDate(division[j].cartel)
		if !iLast.Equal(jLast) {
			return iLast.After(jLast)
		}

		return division[i].contender.ID < division[j].contender.ID
	})
}

// sortFights orders fights by appeal so the best ones are distributed first.
func sortFights(fights []fight) {
	sort.SliceStable(fights, func(i, j int) bool {
		if fights[i].appeal != fights[j].appeal {
			return fights[i].appeal > fights[j].appeal
		}

		if fights[i].favorite.score != fights[j].favorite.score {
			return fights[i].favorite.score > fights[j].favorite.score
		}

		if fights[i].favorite.contender.ID != fights[j].favorite.contender.ID {
			return fights[i].favorite.contender.ID < fights[j].favorite.contender.ID
		}

		return fights[i].underdog.contender.ID < fights[j].underdog.contender.ID
	})
}

// fightAppeal blends star power (average score) with how close the matchup is.
func fightAppeal(favoriteScore, underdogScore int) int {
	average := (favoriteScore + underdogScore) / 2

	gap := favoriteScore - underdogScore
	if gap < 0 {
		gap = -gap
	}

	return clampScore((average + (100 - gap)) / 2)
}

// distributeWeeks spreads fights round-robin across the season, never putting
// more than maxFightsPerWeek on a single card. Fights beyond the season
// capacity are dropped from the tail (lowest appeal first).
func distributeWeeks(fights []fight, weeks int) [][]fight {
	if weeks <= 0 {
		weeks = 1
	}

	capacity := weeks * maxFightsPerWeek
	if len(fights) > capacity {
		fights = fights[:capacity]
	}

	cards := make([][]fight, weeks)
	cursor := 0

	for _, fight := range fights {
		for len(cards[cursor]) >= maxFightsPerWeek {
			cursor = (cursor + 1) % weeks
		}

		cards[cursor] = append(cards[cursor], fight)
		cursor = (cursor + 1) % weeks
	}

	return cards
}

// orderCard splits the appeal between the opening and the main event: the best
// fight closes the card, the second best opens it.
func orderCard(card []fight) []fight {
	sorted := make([]fight, len(card))
	copy(sorted, card)

	sortFights(sorted)

	if len(sorted) <= 1 {
		return sorted
	}

	ordered := make([]fight, len(sorted))
	ordered[0] = sorted[1]

	for i := 2; i < len(sorted); i++ {
		ordered[i-1] = sorted[i]
	}

	ordered[len(sorted)-1] = sorted[0]

	return ordered
}

func totalFights(cartel *model.Cartel) int {
	if cartel == nil {
		return 0
	}

	return cartel.TotalFights
}

func lastFightDate(cartel *model.Cartel) time.Time {
	if cartel == nil || cartel.LastFightDate == nil {
		return time.Time{}
	}

	return *cartel.LastFightDate
}

func clampScore(score int) int {
	if score < 0 {
		return 0
	}

	if score > 100 {
		return 100
	}

	return score
}

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
