DB_URL="postgres://postgres:yourpassword@localhost:5432/chat_db?sslmode=disable"
ATLAS_DEV_URL="postgres://postgres:yourpassword@localhost:5433/atlas_db?sslmode=disable"
MIGRATIONS_DIR="schemas"

setup: deps
	cp .envrc.sample .envrc
	./scripts/install-atlas.sh
	./scripts/build-tools.sh

deps:
	go mod tidy

gen: buf-generate 

compose-up:
	docker compose up -d

run-connect:
	go run cmd/connect

# xo
gen-xo:
	bin/xo schema $(DB_URL) --out ./internal/domain/model

# buf
buf-lint: setup
	./bin/buf lint ./api/proto-spec

buf-generate:
	./bin/buf generate ./api/proto-spec

# migrate
plan-dev:
	docker run --rm -v $(shell pwd)/$(MIGRATIONS_DIR):/migrations --network host arigaio/atlas schema diff --from $(DB_URL) --to file:///migrations/psql.sql --dev-url $(ATLAS_DEV_URL) --format '{{ sql . "  " }}'

migrate-dev:
	docker run --rm -v $(shell pwd)/$(MIGRATIONS_DIR):/migrations --network host arigaio/atlas schema apply --url $(DB_URL) --to file:///migrations/psql.sql --dev-url $(ATLAS_DEV_URL) --auto-approve

# ghz
ghz-run:



# help
help-psql:
	@echo "docker exec -it chat_db psql -U postgres -d chat_db"

# Makefile config
#===============================================================
.SILENT: help

.PHONY: $(shell grep -E -o '^(\._)?[a-z_-]+:' $(MAKEFILE_LIST) | sed 's/://')
