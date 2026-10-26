module github.com/charmingruby/crcs/services/api

go 1.27.0

require (
	crcs/core v0.0.0
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/go-chi/chi/v5 v5.2.5
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.54.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/lib/pq v1.12.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
)

replace crcs/core => ../../pkg/core
