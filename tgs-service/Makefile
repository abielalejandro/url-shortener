HTTP_API_TYPE=http
RPC_API_TYPE=rpc

all: run

build:
	@go build -o bin ./...

run: build
	export API_TYPE=$(HTTP_API_TYPE); ./bin/app

run_grpc: build
	export API_TYPE=$(RPC_API_TYPE);./bin/app

test:
	@go test -v ./test/...

grpc_build_proto:
	protoc --go_out=./api --go_opt=paths=source_relative --go-grpc_out=./api --go-grpc_opt=paths=source_relative tgs.proto
