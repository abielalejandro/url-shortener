
# Url shortener service
It handle the creation and searching os short urls

## How to run
**This service depends of token generator service**

**Read the token generator service documentation first**

**Requirements:**

Before run this service, you must to create the Cassandra db schema (internals/storage/cassandra.schema.cql)

## How to run

### Download dependencies
```
go mod download
```


### Http api
```
make run
```
or

```
go build -o bin ./...
./bin/app
```


### Example calls

#### Create a short token

```
curl --location 'http://localhost:8080/api/v1/short' \
--header 'Content-Type: application/json' \
--data '{
    "url": "http://www.google.com"
}'
```

#### Create a short token with grpc
```
grpcurl -d '{"url": "http://www.google.com"}' -plaintext localhost:9000 api.ShortenerService/Create

```

#### Search the url by token
```
curl --location 'http://localhost:8080/api/v1/short/REPLACE_WIT_TOKEN'

```

#### Search the url by token with grpc
```
grpcurl -d '{"url": "REPLACE_WIT_TOKEN"}' -plaintext localhost:9000 api.ShortenerService/Search

```

#### Health

```
curl --location 'http://localhost:8080/health'
```

#### Health with grpc
```
grpcurl -plaintext localhost:9000 api.ShortenerService/Health

```


### Tech stack

- Go
- Redis
  - Rate limiter
  - Bloom filter
  - Temporal cache
- Cassandra