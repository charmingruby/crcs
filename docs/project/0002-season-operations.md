# ADR-0002: Season operations

- Status: accepted
- Date: 2026-10-27

## Context

After a season exists and its roster is ingested, an operator still needs to
decide the season, contract or reject contenders, and record the outcome of each
scheduled fight. These are three commands, split across the domains that own the
data: seasons and matches belong to `matchmaking`; contenders belong to
`roster`.

## Decision

- `PATCH /v1/seasons/{id}` updates `name`, `status`, and `start_at` partially —
  omitted fields keep their value, an empty body is a validation error, and
  `weeks` is immutable. `status` accepts `accepted` and `rejected`; both are
  reachable from `draft` and reversible between themselves, never back to
  `draft`. Renaming to another season's name is a conflict.
- `PATCH /v1/seasons/{id}/contenders/{cid}` is owned by `roster`, which is the
  domain that holds contenders. It validates the season through the existing
  `SeasonReader`, and answers `404` when the contender does not belong to that
  season. `contracted` and `rejected` are reversible; a status never returns to
  `idle`.
- `POST /v1/seasons/{id}/match/{mid}/result` takes the full result —
  `outcome`, `winner_id`, `method`, `fight_duration`, `scheduled_rounds`,
  `actual_rounds` — and stores it in `season_match_results`, keyed by
  `match_id`. A win (`W`) requires `winner_id`, which must be one of the two
  contenders of the match; a draw (`D`) or a no-contest (`NC`) is stored with no
  winner. A match accepts a single result (`409` afterwards).
- Updates run through a transaction manager in each domain, since all three are
  read-modify-write flows.
- Recording a result does not touch the contender cartel: the cartel remains the
  worker's ingest responsibility, so a season result only affects scoring after
  it is re-ingested.

## Consequences

- The API keeps one writer per aggregate: `matchmaking` writes seasons, matches,
  and results; `roster` writes contenders.
- Season results and the next matchmaking run are decoupled: until the cartel is
  refreshed, the schedule ignores the result.
- A result always carries its outcome, so draws and no-contests stay
  distinguishable from a win while sharing the same table. The cartel is left
  untouched until recording a result is wired to it.
