# ADR-0001: Season matchmaking

- Status: accepted
- Date: 2026-10-27

## Context

`POST /v1/seasons/{id}/matchmaking` must turn a season roster into a fight
schedule. The API checks that the season exists and enqueues an event; the
worker owns the calculations and the persistence.

## Decision

- One event per request: `seasons.matchmaking_requested.v1` with the `SeasonID`.
- Only contenders with `active` status enter the schedule. Contenders without a
  `gender` are skipped, and a fight only happens inside the same
  `weight_class` + `gender`.
- Within a division, contenders are sorted by score and paired adjacently
  (1st vs 2nd, 3rd vs 4th…), so the closest matchups happen first. The higher
  score is the favorite. An odd division leaves its last (lowest score)
  contender without a fight.
- Each contender fights at most once in the whole season, so a contender may
  end up with no fight at all.
- A week card holds at most 5 fights and has no minimum. Fights are spread
  round-robin across the season; when the schedule exceeds `weeks × 5`, the
  excess is dropped from the lowest appeal up. The task still completes
  successfully, even when some weeks end up empty.
- Appeal of a fight blends star power and closeness:
  `(avg(score) + (100 - scoreGap)) / 2`.
- Card order splits the appeal between the ends: the best fight is the main
  event (last position), the second best opens the card (first position); the
  remaining fights stay in descending appeal between them.
- Contender score is `45% record + 35% avg performance rating + 20% experience`,
  clamped to 0–100, where `record = (wins + 0.5 * draws) / totalFights`,
  `experience = min(totalFights, 10) * 10`. A missing cartel scores as a
  newcomer (`record = 0`, `quality = 50`), and a missing rating defaults
  `quality` to 50. Ties break on total fights, last fight date, then ID.
- Regeneration preserves matches with results, including their IDs, weeks, and
  positions. Their contenders are excluded from pairing, and their occupied
  slots count toward the five-fight weekly limit. Only matches without results
  are replaced, in the same transaction.
- Matchmaking locks the season until commit. Result recording takes the same
  lock, so it cannot race with schedule replacement.
- Redelivery of an event is ignored after its first committed processing. A new
  matchmaking request is a new event and intentionally regenerates the schedule.

## Consequences

- A season needs enough active contenders to fill its cards; with a small
  roster later weeks end up with fewer fights, and a roster smaller than the
  number of weeks leaves empty cards.
- Changing the scoring weights or the card rule changes existing schedules on
  the next matchmaking run.
- IDs of matches without results change on each new matchmaking request, even
  when the schedule itself is identical. Redelivery does not change IDs.
- Results already recorded are never deleted. Regeneration cannot reorder
  preserved matches to optimize the appeal of their cards.
