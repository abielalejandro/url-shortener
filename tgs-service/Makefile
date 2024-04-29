HTTP_API_TYPE=http
RPC_API_TYPE=grpc

all: run

build:
	@go build -o bin ./...

run: build
	export API_TYPE=$(HTTP_API_TYPE); ./bin/app

run_grpc: build
	export API_TYPE=$(RPC_API_TYPE);./bin/app

test:
	@go test -v ./test/...