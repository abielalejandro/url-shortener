FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/app .

FROM golang:1.22-alpine3.19
WORKDIR /
RUN apk add --update gettext curl
COPY --from=builder /app/bin/app /app
COPY --from=builder /app/config.yml ./config.yml
COPY --from=builder /app/config.yml ./config.yml
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
EXPOSE 8080

CMD ["/app"]