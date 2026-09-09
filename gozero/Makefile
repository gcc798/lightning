ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
APPLICATION_DIR := $(ROOT_DIR)/application

SYS_API_DIR := $(APPLICATION_DIR)/sys-api
SYS_RPC_DIR := $(APPLICATION_DIR)/sys-rpc

BIN_DIR := $(ROOT_DIR)/bin
STYLE ?= goZero

PG_URL ?= postgres://postgres:post123@127.0.0.1:5433/lightning?sslmode=disable
PG_SCHEMA ?= public
PG_TABLES ?= s_*

.PHONY: help init-tree gen-sys-api gen-sys-rpc gen-sys-model-pg gen-all build-sys-api build-sys-rpc build-all run-sys-api run-sys-rpc package clean

help: ## Display help
	@grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

init-tree: ## Initialize sys-api/sys-rpc directory tree
	@mkdir -p $(SYS_API_DIR)/api
	@mkdir -p $(SYS_API_DIR)/etc
	@mkdir -p $(SYS_API_DIR)/internal/{config,handler,logic,middleware,svc,types}
	@mkdir -p $(SYS_RPC_DIR)/pb
	@mkdir -p $(SYS_RPC_DIR)/client/sysservice
	@mkdir -p $(SYS_RPC_DIR)/etc
	@mkdir -p $(SYS_RPC_DIR)/internal/{config,logic,model,server,svc}
	@mkdir -p $(BIN_DIR)
	@echo "Initialized application tree under $(APPLICATION_DIR)"

gen-sys-api: ## Generate sys-api service code from sys.api
	@echo "Generating sys-api..."
	@goctl api go \
		--api=$(SYS_API_DIR)/sys.api \
		--dir=$(SYS_API_DIR) \
		--style=$(STYLE)

gen-sys-rpc: ## Generate sys-rpc service code from sys.proto
	@echo "Generating sys-rpc..."
	@rm -rf $(SYS_RPC_DIR)/pb/*
	@goctl rpc protoc \
		-I=$(SYS_RPC_DIR) \
		--go_out=$(SYS_RPC_DIR)/pb \
		--go_opt=module=github.com/gcc798/lightning/application/sys-rpc/pb \
		--go-grpc_out=$(SYS_RPC_DIR)/pb \
		--go-grpc_opt=module=github.com/gcc798/lightning/application/sys-rpc/pb \
		--zrpc_out=$(SYS_RPC_DIR) \
		--style=$(STYLE) \
		-m \
		$(SYS_RPC_DIR)/sys.proto

gen-sys-model-pg: ## Generate sqlx models from PostgreSQL datasource (requires PG_URL)
	@if [ -z "$(PG_URL)" ]; then \
		echo "PG_URL is required. Example:"; \
		echo "make gen-sys-model-pg PG_URL='postgres://user:pass@127.0.0.1:5432/db?sslmode=disable'"; \
		exit 1; \
	fi
	@echo "Generating sqlx models from PostgreSQL..."
	@goctl model pg datasource \
		--url="$(PG_URL)" \
		--schema="$(PG_SCHEMA)" \
		--table "$(PG_TABLES)" \
		--dir $(SYS_RPC_DIR)/internal/model \
		--strict

gen-all: gen-sys-api gen-sys-rpc ## Generate api and rpc code

build-sys-api: ## Build sys-api binary
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/sys-api $(SYS_API_DIR)/sys.go

build-sys-rpc: ## Build sys-rpc binary
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/sys-rpc $(SYS_RPC_DIR)/sys.go

build-all: build-sys-api build-sys-rpc ## Build all binaries

run-sys-api: ## Run sys-api service
	@go run $(SYS_API_DIR)/sys.go -f $(SYS_API_DIR)/etc/sys-api.yaml

run-sys-rpc: ## Run sys-rpc service
	@go run $(SYS_RPC_DIR)/sys.go -f $(SYS_RPC_DIR)/etc/sys-rpc.yaml

package: build-all ## Package binaries
	@tar -czf $(BIN_DIR)/sys-services.tar.gz -C $(BIN_DIR) sys-api sys-rpc
	@echo "Package created: $(BIN_DIR)/sys-services.tar.gz"

clean: ## Clean build artifacts
	@rm -rf $(BIN_DIR)
	@echo "Cleaned $(BIN_DIR)"
