run_server:
	go run cmd/bot/main.go cmd/bot/server.go
run_client:
	go run client/client.go
gen:
	protoc --go_out=pkg --go_opt=paths=source_relative --go-grpc_out=pkg --go-grpc_opt=paths=source_relative api/api.proto