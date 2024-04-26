# Url Shortener

### Micro services

<dl>
  <dt>Token generator service (tgs-service) </dt>
  <dd>It handle the sequence for the shorts tokens </dd>
  <dt>Url shortener (shortener-service)</dt>
  <dd>It handle the creation and searching os short urls</dd>
  <dt>Web (shortener-web)</dt>
  <dd>Coming soon a beauty....</dd>  
</dl>

### How to run

**Read the shortener service and token generator service documentation first**

```
cd docker
docker compose docker-compose up -d
```

### Using the web

Access to http://localhost:8080

### Example calls

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

### Tech stack

- Go
- Redis
  - Rate limiter
  - Bloom filter
  - Incremental
  - Temporal cache
- Cassandra
- Consul
  - Service discovery

