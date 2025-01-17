# Makefile to configure commands for all services

# Path to the YAML configuration file
OWNER_CONFIG_FILE=owner-service/config-sample.yml
PRODUCT_CONFIG_FILE=product-service/config-sample.yml
ORDER_CONFIG_FILE=order-service/config-sample.yml

# Variables extracted from YAML using yq

OWNER_DB_NAME := $(shell yq '.database.name' $(OWNER_CONFIG_FILE))
OWNER_DB_USER := $(shell yq '.database.user' $(OWNER_CONFIG_FILE))
OWNER_DB_PASSWORD := $(shell yq '.database.password' $(OWNER_CONFIG_FILE))

PRODUCT_DB_NAME := $(shell yq '.database.name' $(PRODUCT_CONFIG_FILE))
PRODUCT_DB_USER := $(shell yq '.database.user' $(PRODUCT_CONFIG_FILE))
PRODUCT_DB_PASSWORD := $(shell yq '.database.password' $(PRODUCT_CONFIG_FILE))

ORDER_DB_NAME := $(shell yq '.database.name' $(ORDER_CONFIG_FILE))
ORDER_DB_USER := $(shell yq '.database.user' $(ORDER_CONFIG_FILE))
ORDER_DB_PASSWORD := $(shell yq '.database.password' $(ORDER_CONFIG_FILE))

RABBITMQ_USER := $(shell yq '.rabbitmq.user' $(OWNER_CONFIG_FILE))
RABBITMQ_PASSWORD := $(shell yq '.rabbitmq.password' $(OWNER_CONFIG_FILE))

service-build: ## Rebuild image and containers
	DB_PASSWORD=$(OWNER_DB_PASSWORD) DB_NAME=$(OWNER_DB_NAME) DB_USER=$(OWNER_DB_USER) docker-compose -f owner-service/docker-compose.yml up --build -d
	DB_PASSWORD=$(PRODUCT_DB_PASSWORD) DB_NAME=$(PRODUCT_DB_NAME) DB_USER=$(PRODUCT_DB_USER) docker-compose -f product-service/docker-compose.yml up --build -d
	DB_PASSWORD=$(ORDER_DB_PASSWORD) DB_NAME=$(ORDER_DB_NAME) DB_USER=$(ORDER_DB_USER) docker-compose -f order-service/docker-compose.yml up --build -d
	RABBITMQ_USER=$(RABBITMQ_USER) RABBITMQ_PASSWORD=$(RABBITMQ_PASSWORD) docker-compose up --build -d

service-up: ## Start docker services
	DB_PASSWORD=$(OWNER_DB_PASSWORD) DB_NAME=$(OWNER_DB_NAME) DB_USER=$(OWNER_DB_USER) docker compose -f owner-service/docker-compose.yml up -d
	DB_PASSWORD=$(PRODUCT_DB_PASSWORD) DB_NAME=$(PRODUCT_DB_NAME) DB_USER=$(PRODUCT_DB_USER) docker compose -f product-service/docker-compose.yml up -d
	DB_PASSWORD=$(ORDER_DB_PASSWORD) DB_NAME=$(ORDER_DB_NAME) DB_USER=$(ORDER_DB_USER) docker compose -f order-service/docker-compose.yml up -d
	RABBITMQ_USER=$(RABBITMQ_USER) RABBITMQ_PASSWORD=$(RABBITMQ_PASSWORD) docker-compose up -d

service-down: ## Stop services
	docker-compose -f owner-service/docker-compose.yml down
	docker-compose -f product-service/docker-compose.yml down
	docker-compose -f order-service/docker-compose.yml down
	docker-compose down

service-down-add: ## Stop services, volumes and networks
	docker-compose -f owner-service/docker-compose.yml down -v
	docker-compose -f product-service/docker-compose.yml down -v
	docker-compose -f order-service/docker-compose.yml down -v
	docker-compose down -v
