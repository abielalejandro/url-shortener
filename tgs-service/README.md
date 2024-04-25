
# Token generation service
It handle the sequence for the shorts tokens

It has 2 communication apis: http rest (default) and grpc

## How to run

### Http api
```
make run
```
or

```
go build -o bin ./...
./bin/app
```

### Grpc api
```
export API_TYPE=grpc
make run_grpc
```
or

```
export API_TYPE=grpc
go build -o bin ./...
./bin/app
```

### Example calls

#### Create a short token

```
curl --location 'http://localhost:8080/api/v1/next'
```

#### Create a short token with grpc
```
grpcurl -plaintext localhost:9000 api.TgsService/Next

```

#### Health

```
curl --location 'http://localhost:8080/health'
```

#### Health with grpc
```
grpcurl -plaintext localhost:9000 api.TgsService/Health

```

### Tech stack

- Go
- Redis
  - Incremental
- Base62
