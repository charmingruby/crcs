module github.com/charmingruby/crcs/services/outbox

go 1.27.0

require (
	crcs/core v0.0.0
	crcs/platform v0.0.0
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/lib/pq v1.12.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.51 // indirect
)

replace crcs/platform => ../../pkg/platform

replace crcs/core => ../../pkg/core
