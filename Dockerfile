# syntax=docker/dockerfile:1

FROM golang:1.23.3 AS build-stage

# Set destination for COPY
WORKDIR /app

COPY go.mod go.sum embed.go ./
RUN go mod download

COPY ./cmd/sso/main.go ./cmd/sso/main.go
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./migrations ./migrations
COPY ./static ./static

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/ cmd/sso/main.go

# Deploy the application binary into a lean image
FROM ubuntu:22.04 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/bin/main /app/main

EXPOSE 8080
EXPOSE 44844

ENTRYPOINT [ "./app/main" ] 