run_server:
	go run cmd/bot/main.go cmd/bot/server.go
run_client:
	go run client/client.go
gen:
	protoc --go_out=pkg --go_opt=paths=source_relative --plugin=protoc-gen-go=bin/protoc-gen-go \
		   --go-grpc_out=pkg --go-grpc_opt=paths=source_relative --plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		   --grpc-gateway_out=pkg --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true --grpc-gateway_opt=allow_delete_body=true --plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
		   api/api.proto

LOCAL_BIN:=$(CURDIR)/bin
.PHONY: .deps
.deps:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go && \
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
