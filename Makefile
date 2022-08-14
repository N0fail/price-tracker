.PHONY: run_server
run_server:
	go run cmd/bot/main.go cmd/bot/server.go

.PHONY: run_server_cache
run_server_cache:
	go run cmd/bot/main.go cmd/bot/server.go -cache

.PHONY: run_client
run_client:
	go run client/client.go

.PHONY: gen
gen:
	protoc --go_out=pkg --go_opt=paths=source_relative --plugin=protoc-gen-go=bin/protoc-gen-go \
		   --go-grpc_out=pkg --go-grpc_opt=paths=source_relative --plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		   --grpc-gateway_out=pkg --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true --grpc-gateway_opt=allow_delete_body=true --plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
		   --openapiv2_out=swagger --openapiv2_opt=logtostderr=true --plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
		   api/api.proto

LOCAL_BIN:=$(CURDIR)/bin
.PHONY: .deps
.deps:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: up_db
up_db:
	docker-compose build
	docker-compose up -d postgres

.PHONE: down_db
down_db:
	docker-compose down

MIGRATIONS_DIR=./migrations
.PHONY: migration
migration:
	goose -dir=${MIGRATIONS_DIR} create $(NAME) sql

GENERATIONS_DIR=./migrations/generations
.PHONY: generation
generation:
	goose -dir=${GENERATIONS_DIR} create $(NAME) sql
