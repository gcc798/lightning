KRATOS_ROOT := $(CURDIR)
API_DIR := $(KRATOS_ROOT)/api
THIRD_PARTY_DIR := $(KRATOS_ROOT)/third_party
SYS_API_DIR := $(KRATOS_ROOT)/application/sys-api
SYS_RPC_DIR := $(KRATOS_ROOT)/application/sys-rpc
SYS_API_SERVICE_GEN_DIR := $(SYS_API_DIR)/internal/servicegen
SYS_RPC_SERVICE_GEN_DIR := $(SYS_RPC_DIR)/internal/servicegen
PROTO_FILES := $(shell cd $(API_DIR) && find . -name '*.proto' | sed 's|^\./||' | sort)
CONF_PROTO_FILES := application/sys-api/internal/conf/conf.proto application/sys-rpc/internal/conf/conf.proto
PROTOBUF_INCLUDE := /opt/homebrew/include/google/protobuf

export PATH := /Users/guoc/dev/code_go/bin:$(PATH)
export GOCACHE := /tmp/microservice-kit-kratos-go-build

.PHONY: help init prepare-third-party clean-proto clean-conf conf proto openapi proto-all proto-add-health proto-client proto-server-sys-api proto-server-sys-rpc ent wire fmt build-sys-api build-sys-rpc build-all test

help:
	@echo "targets: init prepare-third-party clean-proto clean-conf conf proto openapi proto-all proto-add-health proto-client proto-server-sys-api proto-server-sys-rpc ent wire fmt build-sys-api build-sys-rpc build-all test"

init: prepare-third-party
	go mod tidy

prepare-third-party:
	mkdir -p $(THIRD_PARTY_DIR)/google/api $(THIRD_PARTY_DIR)/google/protobuf $(THIRD_PARTY_DIR)/validate
	cp /Users/guoc/dev/code_go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis/google/api/annotations.proto $(THIRD_PARTY_DIR)/google/api/
	cp /Users/guoc/dev/code_go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis/google/api/http.proto $(THIRD_PARTY_DIR)/google/api/
	cp $(PROTOBUF_INCLUDE)/descriptor.proto $(THIRD_PARTY_DIR)/google/protobuf/
	cp $(PROTOBUF_INCLUDE)/duration.proto $(THIRD_PARTY_DIR)/google/protobuf/
	cp $(PROTOBUF_INCLUDE)/struct.proto $(THIRD_PARTY_DIR)/google/protobuf/
	cp $(PROTOBUF_INCLUDE)/timestamp.proto $(THIRD_PARTY_DIR)/google/protobuf/
	cp /Users/guoc/dev/code_go/pkg/mod/github.com/envoyproxy/protoc-gen-validate@v1.0.4/validate/validate.proto $(THIRD_PARTY_DIR)/validate/

clean-proto:
	rm -rf $(KRATOS_ROOT)/api/system/v1/*.pb.go $(KRATOS_ROOT)/api/system/v1/*_grpc.pb.go $(KRATOS_ROOT)/api/system/v1/*_http.pb.go $(KRATOS_ROOT)/api/system/v1/*_pb.validate.go $(KRATOS_ROOT)/api/system/v1/*_client.go $(KRATOS_ROOT)/system
	find $(API_DIR) -name 'openapi.yaml' -delete
	find $(API_DIR) -name '*.openapi.yaml' -delete
	rm -f $(KRATOS_ROOT)/openapi.yaml

clean-conf:
	rm -f $(SYS_API_DIR)/internal/conf/conf.pb.go $(SYS_RPC_DIR)/internal/conf/conf.pb.go

conf: clean-conf
	protoc \
		-I . \
		--go_out=paths=source_relative:. \
		$(CONF_PROTO_FILES)

proto-add-health:
	cd $(API_DIR) && kratos proto add system/v1/health.proto

proto:
	cd $(API_DIR) && protoc \
		-I . \
		-I $(THIRD_PARTY_DIR) \
		--go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--go-http_out=paths=source_relative:. \
		--go-errors_out=paths=source_relative:. \
		--validate_out=lang=go,paths=source_relative:. \
		$(PROTO_FILES)

openapi:
	tmp_root=$$(mktemp -d /tmp/microservice-kit-openapi.XXXXXX); \
	rm -f $(KRATOS_ROOT)/openapi.yaml; \
	find $(API_DIR) -name 'openapi.yaml' -delete; \
	find $(API_DIR) -name '*.openapi.yaml' -delete; \
	for dir in $$(printf '%s\n' $(PROTO_FILES) | xargs -n1 dirname | sort -u); do \
		inputs=""; \
		index=0; \
		cd $(API_DIR) && for file in $$(find $$dir -maxdepth 1 -name '*.proto' | sort); do \
			index=$$((index + 1)); \
			outdir=$$tmp_root/$$dir/$$index; \
			mkdir -p $$outdir; \
			protoc -I . -I $(THIRD_PARTY_DIR) --openapi_out=fq_schema_naming=true,default_response=false:$$outdir $$file; \
			inputs="$$inputs $$outdir/openapi.yaml"; \
		done; \
		service=$${dir%%/*}; \
		first=$$(printf '%s' "$$service" | cut -c1 | tr '[:lower:]' '[:upper:]'); \
		rest=$$(printf '%s' "$$service" | cut -c2-); \
		title="$$first$$rest API"; \
		version=$$(printf '%s' "$$dir" | awk -F/ '{print $$2}'); \
		if [ -n "$$inputs" ]; then \
			cd $(KRATOS_ROOT) && GOCACHE=$(GOCACHE) go run ./cmd/openapi-merge -title "$$title" -out "$(API_DIR)/$$service/$$version.openapi.yaml" $$inputs; \
		fi; \
	done; \
	rm -rf $$tmp_root

proto-all: conf clean-proto prepare-third-party proto openapi proto-client proto-server-sys-api proto-server-sys-rpc

proto-client:
	cd $(API_DIR) && for file in $(PROTO_FILES); do \
		kratos proto client $$file -p $(THIRD_PARTY_DIR); \
	done

proto-server-sys-api:
	cd $(API_DIR) && for file in $(PROTO_FILES); do \
		mkdir -p $(SYS_API_SERVICE_GEN_DIR); \
		kratos proto server $$file -t $(SYS_API_SERVICE_GEN_DIR); \
	done

proto-server-sys-rpc:
	cd $(API_DIR) && for file in $(PROTO_FILES); do \
		mkdir -p $(SYS_RPC_SERVICE_GEN_DIR); \
		kratos proto server $$file -t $(SYS_RPC_SERVICE_GEN_DIR); \
	done

ent:
	ent generate ./application/sys-rpc/ent/schema

wire:
	GOSUMDB=off go run github.com/google/wire/cmd/wire ./application/sys-api/cmd/server ./application/sys-rpc/cmd/server

fmt:
	go fmt ./...

build-sys-api:
	go build -o $(KRATOS_ROOT)/bin/sys-api ./application/sys-api/cmd/server

build-sys-rpc:
	go build -o $(KRATOS_ROOT)/bin/sys-rpc ./application/sys-rpc/cmd/server

build-all: build-sys-api build-sys-rpc

test:
	go test ./...
