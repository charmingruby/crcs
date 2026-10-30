# Lab

A lean Go foundation for proving ideas fast.

## What is this?

A place to validate Go ideas quickly. Copy a domain, tweak, and see if it works.

## Getting Started

```sh
cp .env.example .env
docker compose up -d
task
```

## Asynchronous processing

Run the API, worker, and outbox service as separate processes after applying migrations:

```sh
task mig:up
go run ./services/api/cmd/api
go run ./services/worker/cmd/worker
task outbox
```

The API queues requests in Postgres even when Kafka is unavailable. The outbox service
publishes pending events, and the worker deduplicates their deliveries.
`services/outbox` has its own Go module, configuration, and executable. It uses
`DATABASE_URL` and `KAFKA_BROKERS`, with the same defaults as the other services.

`task build` builds all modules. `task test` runs tests; asynchronous regression
tests use an isolated schema in `TEST_DATABASE_URL`. Set `TEST_KAFKA_BROKER` to also
exercise real Kafka publication and republication (`task test:async`).

## Docs

- [Architecture overview](docs/internals/architecture.md)
- [Platform](docs/internals/platform.md)
- [Glossary](docs/internals/glossary.md)
- [Coding patterns](docs/internals/coding-patterns.md)
- [Commands reference](docs/internals/commands.md)
- [Messaging](docs/internals/messaging.md)
- [External integrations](docs/internals/external-integrations.md)
- [Cross-module reads](docs/internals/cross-module-reads.md)
- [Versioning](docs/internals/versioning.md)
- [Testing](docs/internals/testing.md)
