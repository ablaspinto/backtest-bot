
# Database
DB_DSN ?= host=localhost user=postgres password=postgres dbname=marketsim sslmode=disable
MIGRATIONS_DIR ?= internal/adapters/postgresql/migrations

# ==================== Dev ====================
.PHONY: dev
dev:
	docker compose up -d
	air

.PHONY: run
run:
	go run ./cmd/*.go

# ==================== Docker ====================
.PHONY: up
up:
	docker compose up -d

.PHONY: down
down:
	docker compose down

.PHONY: logs
logs:
	docker compose logs -f

# ==================== Database ====================
.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

.PHONY: migrate-down
migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

.PHONY: migrate-reset
migrate-reset:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" reset

.PHONY: migrate-status
migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

.PHONY: migrate-create
migrate-create:
	@read -p "Migration name: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql

.PHONY: db-shell
db-shell:
	docker exec -it backtest-postgres psql -U postgres -d marketsim

# ==================== Build ====================
.PHONY: build
build:
	go build -o ./tmp/main ./cmd/.

.PHONY: clean
clean:
	rm -rf ./tmp

# ==================== Help ====================
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make dev            - Start docker + air (hot reload)"
	@echo "  make run            - Run the server"
	@echo "  make up             - Start docker containers"
	@echo "  make down           - Stop docker containers"
	@echo "  make logs           - Tail docker logs"
	@echo "  make migrate-up     - Run all migrations"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make migrate-reset  - Rollback all migrations"
	@echo "  make migrate-status - Show migration status"
	@echo "  make migrate-create - Create new migration"
	@echo "  make db-shell       - Open psql shell"
	@echo "  make build          - Build binary"
