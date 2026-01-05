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