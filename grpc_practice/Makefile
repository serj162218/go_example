protoc:
	protoc --go_out=. --go-grpc_out=. ./proto/order.proto

# Run the server
run-server:
	go run server/main.go

# Run the client
run-client:
	go run client/main.go