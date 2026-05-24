.PHONY: all build build-server build-worker build-cli migrate-up migrate-down generate run-server run-worker run-all dev test lint clean

# Variables
BINARY_SERVER=bin/server
BINARY_WORKER=bin/worker
BINARY_CLI=bin/atp
BINARY_MIGRATE=bin/migrate
GO=go
NPM=npm

all: generate build

## build: Build all binaries
build: build-server build-worker build-cli

## build-server: Build the API server
build-server:
	$(GO) build -o $(BINARY_SERVER) ./cmd/server

## build-worker: Build the worker agent
build-worker:
	$(GO) build -o $(BINARY_WORKER) ./cmd/worker

## build-cli: Build the CLI tool
build-cli:
	$(GO) build -o $(BINARY_CLI) ./cmd/cli

## build-migrate: Build the migration runner
build-migrate:
	$(GO) build -o $(BINARY_MIGRATE) ./cmd/migrate

## run-server: Run the API server
run-server: build-server
	./$(BINARY_SERVER)

## run-worker: Run the worker agent
run-worker: build-worker
	./$(BINARY_WORKER)

## run-all: Run server and worker together (for dev)
run-all: build-server build-worker
	@echo "Starting server and worker..."
	@./$(BINARY_SERVER) & SERVER_PID=$$!; \
	./$(BINARY_WORKER) & WORKER_PID=$$!; \
	trap "kill $$SERVER_PID $$WORKER_PID 2>/dev/null" EXIT; \
	wait

## migrate-up: Run database migrations up
migrate-up:
	$(GO) run ./cmd/migrate -direction up

## migrate-down: Run database migrations down
migrate-down:
	$(GO) run ./cmd/migrate -direction down

## generate: Run sqlc code generation
generate:
	sqlc generate

## dev: Start frontend dev server
dev:
	cd web && $(NPM) run dev

## install-frontend: Install frontend dependencies
install-frontend:
	cd web && $(NPM) install

## test: Run all tests
test:
	$(GO) test ./... -v -count=1

## test-short: Run short tests only
test-short:
	$(GO) test ./... -v -short -count=1

## lint: Run linter
lint:
	golangci-lint run ./...

## clean: Remove build artifacts
clean:
	rm -rf bin/ web/dist/ web/node_modules/

## docker-up: Start all services with Docker Compose
docker-up:
	docker-compose -f deployments/docker-compose.yml up -d

## docker-down: Stop all services
docker-down:
	docker-compose -f deployments/docker-compose.yml down

## docker-build: Build Docker images
docker-build:
	docker-compose -f deployments/docker-compose.yml build

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
