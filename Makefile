# Load .env file
ifneq ("$(wildcard .env)","")
include .env
export $(shell sed 's/=.*//' .env)
endif

APP_VERSION?=0.1.0
IMAGE?=players-service
GOCMD=go
GOBUILD=$(GOCMD) build
CONTAINERTOOL?=docker
COMMIT_HASH?=$(shell git describe --dirty --tags --always)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS?="-X github.com/fernandoocampo/players/internal/application.version=${APP_VERSION} -X github.com/fernandoocampo/players/internal/application.commitHash=${COMMIT_HASH} -X github.com/fernandoocampo/players/internal/application.buildDate=${BUILD_DATE} -s -w"
CONTAINER_COMPOSE_TOOL ?= docker compose
SRC_FOLDER=cmd/playersd
BINARY_NAME=players-service
BINARY_UNIX=$(BINARY_NAME)-linux
BINARY_DARWIN=$(BINARY_NAME)-darwin

.PHONY: all
all: help

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: tidy
tidy: ## Run go mod tidy to organize dependencies.
	@$(GOCMD) mod tidy

.PHONY: run
run: ## Run run application locally.
	@$(GOCMD) run -ldflags ${LDFLAGS} cmd/playersd/main.go

.PHONY: lint
lint: ## Run lint
	@$(CONTAINERTOOL) run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.61.0-alpine golangci-lint run

.PHONY: test
test: ## Run unit tests
	@$(GOCMD) test -count=1 -race ./...

.PHONY: compile-proto
compile-proto: ## compile proto
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    pkg/pb/players/players.proto

.PHONY: build-linux
build-linux: ## Build binary for Linux taking GOARCH from env
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux ${GOBUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_UNIX} ./${SRC_FOLDER}/main.go

.PHONY: build-linux-amd-64
build-linux-amd-64: ## Build binary for Linux amd64
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GOBUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_UNIX}-amd64 ./${SRC_FOLDER}/main.go

.PHONY: build-linux-arm-64
build-linux-arm-64: ## Build binary for Linux amd64
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 ${GOBUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_UNIX}-arm64 ./${SRC_FOLDER}/main.go

.PHONY: build-image
build-image: ## build container image
	@${CONTAINERTOOL} build \
	--build-arg appVersion=${APP_VERSION} \
	--build-arg buildDate=${BUILD_DATE} \
	--build-arg commitHash=${COMMIT_HASH} \
	-f deploy/Dockerfile \
	-t ${IMAGE}:local-${APP_VERSION} .

.PHONY: api-db-up
api-db-up: ## start a containerized database
	@$(CONTAINER_COMPOSE_TOOL) -f deploy/docker-compose.yaml up --build database

.PHONY: api-db-down
api-db-down: ## stop a containerized database
	@$(CONTAINER_COMPOSE_TOOL) -f deploy/docker-compose.yaml down database

.PHONY: api-tracer-up
api-tracer-up: ## start a containerized tracer collector plus jaeger
	@$(CONTAINER_COMPOSE_TOOL) -f deploy/docker-compose.yaml up --build otel-collector

.PHONY: api-tracer-down
api-tracer-down: ## stop a containerized tracer collector plus jaeger
	@$(CONTAINER_COMPOSE_TOOL) -f deploy/docker-compose.yaml down otel-collector

.PHONY: api-up
api-up: ## start application plus db
	@$(CONTAINER_COMPOSE_TOOL) -f deploy/docker-compose.yaml up --build

.PHONY: api-down
api-down: ## stop application plus database
	@$(CONTAINER_COMPOSE_TOOL) -f deploy/docker-compose.yaml down

.PHONY: connect-db
connect-db: ## connect to postgresql database
	psql -U $(PLAYERS_POSTGRES_PLAYER) -h localhost -p $(PLAYERS_POSTGRES_PORT)

.PHONY: migration-up
migration-up: ## run migration up
	@$(CONTAINERTOOL) run --rm -it -v ./migrations:/migrations \
		--network deploy_mynetwork migrate/migrate \
		-path=/migrations/ \
		-database "postgres://$(PLAYERS_POSTGRES_PLAYER):$(PLAYERS_POSTGRES_PASSWORD)@database:$(PLAYERS_POSTGRES_PORT)/$(PLAYERS_POSTGRES_DB)?connect_timeout=10&sslmode=disable" \
		up 1

.PHONY: migration-down
migration-down: ## Apply all or N down migrations
	@$(CONTAINERTOOL) run --rm -it -v ./migrations:/migrations \
		--network deploy_mynetwork migrate/migrate \
		-path=/migrations/ \
		-database "postgres://$(PLAYERS_POSTGRES_PLAYER):$(PLAYERS_POSTGRES_PASSWORD)@database:$(PLAYERS_POSTGRES_PORT)/$(PLAYERS_POSTGRES_DB)?connect_timeout=10&sslmode=disable" \
		down 1

.PHONY: e2e-test-storage
e2e-test-storage: ## Run e2e tests for storage package
	@$(GOCMD) test -v -run ^Test$ \
		github.com/fernandoocampo/players/internal/adapters/storages \
		-e2e-test

.PHONY: e2e-test-pg-connection
e2e-test-pg-connection: ## Run e2e test to check postgress connection
	@$(GOCMD) test -v -run ^TestPostgresqlConnection$ \
		github.com/fernandoocampo/players/internal/adapters/storages \
		-e2e-test

.PHONY: e2e-test-save-db
e2e-test-save-db: ## Run e2e test to save player
	@$(GOCMD) test -v -run ^TestSavePlayer$ \
		github.com/fernandoocampo/players/internal/adapters/storages \
		-e2e-test

.PHONY: e2e-test-update-db
e2e-test-update-db: ## Run e2e test to update player
	@$(GOCMD) test -v -run ^TestUpdatePlayer$ \
		github.com/fernandoocampo/players/internal/adapters/storages \
		-e2e-test

.PHONY: e2e-test-delete-db
e2e-test-delete-db: ## Run e2e test to delete player
	@$(GOCMD) test -v -run ^TestDeletePlayer$ \
		github.com/fernandoocampo/players/internal/adapters/storages \
		-e2e-test

.PHONY: e2e-test-get-db
e2e-test-get-db: ## Run e2e test to get player by id
	@$(GOCMD) test -v -run ^TestGetPlayerByID$ \
		github.com/fernandoocampo/players/internal/adapters/storages \
		-e2e-test

.PHONY: e2e-test-email-nickname-db
e2e-test-email-nickname-db: ## Run e2e test to check if players with given email and nickname exists
	@$(GOCMD) test -v -run ^TestGetPlayersByNicknameOrEmail$ \
		github.com/fernandoocampo/players/internal/adapters/storages \
		-e2e-test

.PHONY: e2e-test-email-nickname-exist-db
e2e-test-email-nickname-exist-db: ## Run e2e test to check if players with given email and nickname exists
	@$(GOCMD) test -v -run ^TestGetPlayersByNicknameOrEmailBothExist$ \
		github.com/fernandoocampo/players/internal/adapters/storages \
		-e2e-test

.PHONY: e2e-test-email-nickname-ignore-db
e2e-test-email-nickname-ignore-db: ## Run e2e test to check if players with given email and nickname exists
	@$(GOCMD) test -v -run ^TestGetPlayersByNicknameOrEmailIgnore$ \
		github.com/fernandoocampo/players/internal/adapters/storages \
		-e2e-test

.PHONY: e2e-test-search-db
e2e-test-search-db: ## Run e2e test to search players
	@$(GOCMD) test -v -run ^TestSearchPlayers$ \
		github.com/fernandoocampo/players/internal/adapters/storages \
		-e2e-test

.PHONY: e2e-test-grpc-create
e2e-test-grpc-create: ## Run e2e test to save a player using grpc endpoint
	@$(GOCMD) test -v -run ^TestE2ECreatePlayer$ \
		github.com/fernandoocampo/players/internal/adapters/grpc \
		-e2e-test

.PHONY: e2e-test-grpc-update
e2e-test-grpc-update: ## Run e2e test to update a player using grpc endpoint
	@$(GOCMD) test -v -run ^TestE2EUpdatePlayer$ \
		github.com/fernandoocampo/players/internal/adapters/grpc \
		-e2e-test

.PHONY: e2e-test-grpc-delete
e2e-test-grpc-delete: ## Run e2e test to delete a player using grpc endpoint
	@$(GOCMD) test -v -run ^TestE2EDeletePlayer$ \
		github.com/fernandoocampo/players/internal/adapters/grpc \
		-e2e-test

.PHONY: e2e-test-grpc-search
e2e-test-grpc-search: ## Run e2e test to search players using grpc endpoint
	@$(GOCMD) test -v -run ^TestE2ESearchPlayers$ \
		github.com/fernandoocampo/players/internal/adapters/grpc \
		-e2e-test

.PHONY: e2e-test-grpc
e2e-test-grpc: ## Run e2e test for all grpc endpoints
	@$(GOCMD) test -v -run ^Test$ \
		github.com/fernandoocampo/players/internal/adapters/grpc \
		-e2e-test
