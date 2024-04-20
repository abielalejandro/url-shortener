FROM golang:1.22 as builder
COPY go.mod go.sum ./
RUN go mod download
COPY -r * ./
RUN CGO_ENABLED=0 GOOS=linux go build -o bin ./...

FROM golang:1.22
WORKDIR /app
ARG GO_APP=bin/app
COPY --from=builder ${GO_APP} .
EXPOSE 8080
USER nonroot:nonroot
CMD ["./app"]