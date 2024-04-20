FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o bin ./...

FROM golang:1.22-alpine3.19
WORKDIR /
COPY --from=builder /app/bin/app /app
COPY --from=builder /app/config /config
EXPOSE 8080
EXPOSE 9000

CMD ["/app"]