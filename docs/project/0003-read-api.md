# ADR-0003: Read API

- Status: accepted
- Date: 2026-10-29

## Context

The API could create seasons, ingest contenders, run matchmaking, and record
results, but a client had no way to read any of it back. Reads span three
domains: `matchmaking` owns seasons, matches, and results; `roster` owns
contenders and cartels; background work is tracked in `tasks`, whose metadata
already carries the `season_id`. Another track owns a future `event_id` column,
so the read API must not depend on it.

## Decision

- Routes, all behind the same bearer auth as the writes:
  - `GET /v1/seasons` — paginated, optional `status` filter.
  - `GET /v1/seasons/{id}` — one season.
  - `GET /v1/seasons/{id}/matches` — the season card, optional `week` filter.
  - `GET /v1/seasons/{id}/matches/{mid}` — one match with its result.
  - `GET /v1/seasons/{id}/contenders` — paginated, optional `status`,
    `weight_class`, and `gender` filters.
  - `GET /v1/seasons/{id}/contenders/{cid}` — one contender with its cartel.
  - `GET /v1/tasks` — paginated, optional `season_id`, `type`, and `status`
    filters.
  - `GET /v1/tasks/{id}` — one task.
- Every read goes through a usecase and a concrete repository; endpoints only
  parse, map, and write. Validation (enums, season existence) stays in the
  usecase and maps to the usual `404` / `400`.
- Lists answer `httpx.Page[T]`: `{items, total, limit, offset}`. `limit`
  defaults to 20 and is clamped to 1..100, `offset` below zero becomes 0.
- The match payload is IDs-only: it carries `favorite_contender_id`,
  `underdog_contender_id`, and the optional result. The client composes a card
  by joining it with the contenders of the same season, so no read duplicates
  contender fields.
- Season matches are not paginated (`limit`/`offset` are `0`, `total` is the
  number of items); a season card is bounded by `weeks × 5`. Matches and their
  results are fetched in two queries and joined in memory, never per match.
- The task domain is read-only. Tasks are correlated to a season through
  `metadata->>'season_id'`, backed by the expression index
  `idx_tasks_metadata_season_id`. There is no `event_id` in the read contract.
- Contender and task lists run list + count as two queries over the same
  filters, so `total` reflects the filter, not the page.

## Consequences

- Reads are eventually consistent with worker writes (task progress, cartel
  updates): the API only reads what is already persisted.
- A new correlation key for tasks needs a new expression index; filtering on
  unindexed metadata stays a scan.
- The IDs-only match payload keeps one source of truth for contender data, but
  clients must fetch the roster to render names.
- Adding a field to a read DTO is a versioned change in `http/endpoint`, never
  a change to the domain model.
