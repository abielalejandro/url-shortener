all: run

build:
	@go build -o bin ./...

run: build
	export API_TYPE=$(HTTP_API_TYPE); ./bin/app

test:
	@go test -v ./test/...