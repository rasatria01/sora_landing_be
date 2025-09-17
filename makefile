# ==========
# Variables
# ==========
DB_USER = sora
DB_PASS = 12345678
DB_NAME = sora_db
DB_HOST = localhost
DB_PORT = 5433
MIGRATIONS_DIR = ./migrations
DATABASE_URL = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# ==========
# Commands
# ==========
.PHONY: help install-migrate createdb dropdb migrateup migratedown newmigration run test

help:
	@echo "Makefile commands:"
	@echo "  make install-migrate     - installs migrate CLI if missing"
	@echo "  make createdb            - Create database (requires createdb client)"
	@echo "  make dropdb              - Drop database"
	@echo "  make migrateup           - Apply all migrations"
	@echo "  make migratedown         - Rollback last migration"
	@echo "  make newmigration n=NAME - Create new migration file"
	@echo "  make run                 - Run the API server (go run)"
	@echo "  make test                - Run go tests"

install-migrate:
	@which migrate >/dev/null 2>&1 || ( \
	curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz \
	  | tar xvz && sudo mv migrate /usr/local/bin/ )

createdb:
	@echo "Creating database..."
	PGPASSWORD=$(DB_PASS) createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME) || true

dropdb:
	@echo "Dropping database..."
	PGPASSWORD=$(DB_PASS) dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) --if-exists $(DB_NAME)

migrateup:
	@echo "Running migrations up..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

migratedown:
	@echo "Rolling back last migration..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

newmigration:
ifndef n
	$(error n is not set. Use: make newmigration n=add_users_table)
endif
	@echo "Creating new migration: $(n)"
	migrate create -ext sql -dir $(MIGRATIONS_DIR) $(n)

run:
	@echo "Running API..."
	# set DB URL in env (adjust as needed)
	@export DATABASE_URL=$(DATABASE_URL) && go run .

test:
	@echo "Running tests..."
	go test ./...
