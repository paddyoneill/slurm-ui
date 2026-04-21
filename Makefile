.PHONY: run
run:
	@go run main.go

GOOSE_DRIVER ?= sqlite3
GOOSE_DBSTRING ?= ./local.db
GOOSE_MIGRATION_DIR ?= ./internal/db/migrations
GOOSE ?= GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) GOOSE_MIGRATION_DIR=$(GOOSE_MIGRATION_DIR) go tool goose

.PHONY: db-up
db-up:
	@$(GOOSE) up

.PHONY: db-down
db-down:
	@$(GOOSE) down

.PHONY: db-reset
db-reset:
	@$(GOOSE) reset

.PHONY: db-status
db-status:
	@$(GOOSE) status
