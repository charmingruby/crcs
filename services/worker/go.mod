module github.com/charmingruby/crcs/services/worker

go 1.27.0

require (
	crcs/core v0.0.0-00010101000000-000000000000
	crcs/event v0.0.0-00010101000000-000000000000
	crcs/platform v0.0.0
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/charmingruby/fsm v0.0.2
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/lib/pq v1.12.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/segmentio/kafka-go v0.4.51 // indirect
)

replace crcs/core => ../../pkg/core

replace crcs/event => ../../pkg/event

replace crcs/platform => ../../pkg/platform
