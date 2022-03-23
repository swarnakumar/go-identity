define load-env
	$(eval ENV_FILE := .env.$(1))
	@echo " - setup env $(ENV_FILE)"
	$(eval include .env.$(1))
	$(eval export)
	$(eval export ENV=$(1))
	$(eval export DB_CONN="user=$(PG_USER) password=$(PG_PASSWORD) host=$(PG_HOST) port=$(PG_PORT) dbname=$(PG_DATABASE) sslmode=disable")
endef

sync-db:
	$(call load-env,dev)
	PGPASSWORD=$(PG_PASSWORD) psql -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -tc "SELECT 1 FROM pg_database WHERE datname = '${PG_DATABASE}'" | grep -q 1 || PGPASSWORD=$(PG_PASSWORD) psql -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -c "CREATE DATABASE ${PG_DATABASE}"
	goose --dir=db/sql/schema postgres $(DB_CONN) up
	goose --dir=db/sql/schema postgres $(DB_CONN) status

setup-test-db:
	$(call load-env,test)
	PGPASSWORD=$(PG_PASSWORD) psql -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -c "DROP DATABASE IF EXISTS ${PG_DATABASE}"
	PGPASSWORD=$(PG_PASSWORD) psql -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -c "CREATE DATABASE ${PG_DATABASE}"
	goose --dir=db/sql/schema postgres $(DB_CONN) up
	goose --dir=db/sql/schema postgres $(DB_CONN) status
	sqlc generate

test: setup-test-db
	go test ./... -v

make-sql:
	sqlc generate

build-dev-css:
	cd web/frontend && npm run build

dev: sync-db build-dev-css
	reflex -c reflex.conf -d fancy -e

build-go: make-sql
	go build -o build/identity

build-css:
	cd web/frontend && npx tailwindcss -i ./input.css -o ../static/main.css --minify

create-user: build-go sync-db make-sql
	cd build && ./identity create-user

build: build-css make-sql
	go build -o build/identity

run: build-go build-css sync-db
	cd build && ./identity run-server

build-http-server: build-go build-css
	go build -o build/server ./cmd/web

build-cli: make-sql
	go build -o build/cli ./cmd/cli
