build:
	@go build -o bin ./...

run: build
	@./bin/app

test:
	@go test -v ./test/...

rpc:
	protoc --go_out=./api --go_opt=paths=source_relative --go-grpc_out=./api --go-grpc_opt=paths=source_relative tgs.proto
