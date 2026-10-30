package repository

import (
	"context"

	"crcs/core/model"
	"crcs/platform/store/postgrex"
)

const (
	deleteSeasonMatchesQuery    = "delete season matches without result"
	createSeasonMatchQuery      = "create season match"
	findMatchesWithResultsQuery = "find matches with results"
)

var matchQueries = map[string]string{
	findMatchesWithResultsQuery: `
		SELECT sm.* FROM season_matches sm
		WHERE sm.season_id=$1 AND sm.deleted_at IS NULL
		AND EXISTS (SELECT 1 FROM season_match_results r
			WHERE r.match_id=sm.id AND r.deleted_at IS NULL)
		ORDER BY sm.week, sm.position`,
	deleteSeasonMatchesQuery: `
		DELETE FROM season_matches sm
		WHERE
			sm.season_id=$1 AND
			NOT EXISTS (
				SELECT 1 FROM season_match_results r
				WHERE r.match_id=sm.id AND r.deleted_at IS NULL
			)`,
	createSeasonMatchQuery: `
		INSERT INTO season_matches
		(id, season_id, favorite_contender_id, underdog_contender_id, week, position, weight_class, gender, created_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
}

type MatchRepository struct {
	stmts *postgrex.Statements
}

func NewMatchRepository(db postgrex.Querier) (*MatchRepository, error) {
	stmts, err := postgrex.Prepare(db, matchQueries)
	if err != nil {
		return nil, err
	}

	return &MatchRepository{stmts: stmts}, nil
}

func (r *MatchRepository) ReplaceBySeason(
	ctx context.Context,
	seasonID string,
	matches []*model.Match,
) error {
	if err := postgrex.Exec(ctx, r.stmts, deleteSeasonMatchesQuery, seasonID); err != nil {
		return err
	}

	for _, match := range matches {
		if err := postgrex.Exec(ctx, r.stmts, createSeasonMatchQuery,
			match.ID,
			match.SeasonID,
			match.FavoriteContenderID,
			match.UnderdogContenderID,
			match.Week,
			match.Position,
			match.WeightClass,
			match.Gender,
			match.CreatedAt,
		); err != nil {
			return err
		}
	}

	return nil
}

func (r *MatchRepository) FindWithResults(ctx context.Context, seasonID string) ([]*model.Match, error) {
	return postgrex.Many[*model.Match](ctx, r.stmts, findMatchesWithResultsQuery, seasonID)
}
