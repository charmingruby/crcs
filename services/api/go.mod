module github.com/charmingruby/crcs/services/api

go 1.27.0

require (
	crcs/core v0.0.0
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/go-chi/chi/v5 v5.2.5
	github.com/go-playground/validator/v10 v10.30.2
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.12.3
	github.com/segmentio/kafka-go v0.4.51
	golang.org/x/crypto v0.54.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/klauspost/compress v1.18.6 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.16 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)

replace crcs/core => ../../pkg/core
