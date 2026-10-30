# ADR-0004: Durable asynchronous delivery

- Status: accepted
- Date: 2026-10-29

## Context

Ingestion and matchmaking requests must survive broker outages, and repeated
delivery must not repeat business writes or create duplicate tasks.

## Decision

- The API's shared outbox publisher (`services/api/internal/shared/outbox`)
  writes to the Postgres `outbox`. A successful response
  means the request is persisted, not that Kafka has received it. These commands
  currently have no accompanying business write. The outbox publisher also accepts
  a transaction when an event needs to be atomic with other writes.
- `services/outbox` is an independent service with its own Go module,
  configuration, and executable (`task outbox`). Its outbox domain owns event
  selection, publication, and acknowledgement through a usecase and concrete
  repository. This logic does not live in the shared messaging infrastructure.
  One instance is sufficient for this PoC. It polls every second, locks a pending row, publishes
  synchronously with broker acknowledgement, then sets `published_at`. Failures
  leave the event pending for the next poll; publication has a 15-second timeout.
- The outbox ID travels in the Kafka `event_id` header, outside the versioned
  event payload. Republishing retains that ID. For older messages without the
  header, the reader uses topic + partition + offset as delivery identity.
- Each worker handler inserts `(listener, message_id)` into `inbox` and processes
  the whole event in that same transaction. Conflicting deliveries wait for its
  outcome: after commit they are skipped; after rollback they can execute.
- Task writes and business effects commit together. Infrastructure failures roll
  everything back, allowing retry without duplicate tasks or partial writes.
- Ingestion keeps per-contender savepoints: an integrity violation rejects only
  that contender. The task's final failures and successful contenders are committed
  together. An infrastructure failure rolls back the entire batch instead.

## Consequences

- Delivery remains at-least-once. A crash after Kafka acknowledges but before the
  relay commits may republish an event; the inbox makes its database effects
  idempotent. This is not an exactly-once Kafka protocol.
- Tasks become visible only when event processing commits; live progress and
  failed infrastructure attempts are not persisted. Whole-event transactions are
  a deliberate tradeoff for small PoC batches.
- A new HTTP request creates a new event, even with the same payload. HTTP request
  deduplication is outside this scope.
- There is no DLQ, retry limit, or automatic retention cleanup. A failing oldest
  event is logged and retried, so it can hold up later events until corrected.
- Listener infrastructure errors still stop the worker without acknowledging the
  message; restarting the worker retries it.
