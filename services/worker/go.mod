module github.com/charmingruby/crcs/services/worker

go 1.27.0

require (
	crcs/platform v0.0.0
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/segmentio/kafka-go v0.4.51 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
)

replace crcs/core => ../../pkg/core

replace crcs/event => ../../pkg/event

replace crcs/platform => ../../pkg/platform
