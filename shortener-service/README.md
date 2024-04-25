
# Url shortener service
It handle the creation and searching os short urls

## How to run
**This service depends of token generator service**

**Read the token generator service documentation first**

**Requirements:**

Before run this service, you must to create the Cassandra db schema (internals/storage/cassandra.schema.cql)

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

#### Create a short token

```
curl --location 'http://localhost:8080/api/v1/short' \
--header 'Content-Type: application/json' \
--data '{
    "url": "http://www.google.com"
}'
```

#### Search the url by token
```
curl --location 'http://localhost:8080/api/v1/short/REPLACE_WIT_TOKEN'

```

#### Health

```
curl --location 'http://localhost:8080/health'
```


### Tech stack

- Go
- Redis
  - Rate limiter
  - Bloom filter
  - Temporal cache
- Cassandra