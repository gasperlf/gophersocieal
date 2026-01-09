#include .envrc
MIGRATE_DIR = ./cmd/migrate/migrations
DATABASE_URL= "postgres://admin:pass@localhost/social?sslmode=disable"

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATE_DIR) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path $(MIGRATE_DIR) -database=$(DATABASE_URL) up 

.PHONY: migrate-down
migrate-down:
	@migrate -path $(MIGRATE_DIR) -database=$(DATABASE_URL) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt

.PHONY: run
run: gen-docs
	@go run cmd/api/*.go