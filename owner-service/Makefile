# Makefile for Dev Environment
# TODO: Find a way to parse configurations in yml

# Path to the YAML configuration file
CONFIG_FILE=config-sample.yml

# Variables extracted from YAML using yq
APP_NAME := $(shell yq '.app.name' $(CONFIG_FILE))
APP_VERSION := $(shell yq '.app.env' $(CONFIG_FILE))

DB_PROTOCOL := $(shell yq '.database.protocol' $(CONFIG_FILE))
DB_HOST := $(shell yq '.database.host' $(CONFIG_FILE))
DB_PORT := $(shell yq '.database.port' $(CONFIG_FILE))
DB_NAME := $(shell yq '.database.name' $(CONFIG_FILE))
DB_USER := $(shell yq '.database.user' $(CONFIG_FILE))
DB_PASSWORD := $(shell yq '.database.password' $(CONFIG_FILE))

DSN := $(DB_PROTOCOL)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)

ARG ?= 

.PHONY: default install service-up service-down db-docs db-create db-drop db-cli \
        migrate-up migrate-down redis-cli dev lint build start swag test sqlc-gen \
		protoc_gen

default: install ## Getting started

install: ## Install dependencies
	go mod download
	go install github.com/air-verse/air@latest
	brew install yq

service-build: ## Rebuild image and containers
	DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME) DB_USER=$(DB_USER) docker-compose up --build -d

service-up: ## Start docker services
	DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME) DB_USER=$(DB_USER) docker-compose up -d

service-down: ## Stop services
	docker-compose down

service-down-add: ## Stop services, volumes and networks
	docker-compose down -v

redis-cli: ## Connect to redis using command line interface
	docker exec -it owner_redis redis-cli

dev: ## Start development server
	air

lint: ## Run linter
	golangci-lint run ./...

print_dsn:
	echo $(DSN)

build: ## Build binary
	go build -o ./bin/$(APP_NAME) ./cmd/http/main.go

start: build ## Start binary
	./bin/$(APP_NAME)

swag: ## Generate swagger documentation
	swag fmt
	swag init -g ./cmd/http/main.go -o ./docs --parseInternal true

test: ## Run tests
	go test -v ./... -race -cover -timeout 30s -count 1 -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

air: ## Start air
	air

protoc-gen:
	protoc internal/core/service/user/*.proto \
                --go_out=internal/core/service/user \
                --go-grpc_out=internal/core/service/user \
                --go_opt=paths=source_relative \
                --go-grpc_opt=paths=source_relative \
                --proto_path=internal/core/service/user