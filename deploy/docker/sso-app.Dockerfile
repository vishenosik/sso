# syntax=docker/dockerfile:1

FROM golang:1.24.2 AS build-stage

# Set destination for COPY
WORKDIR /app

COPY go.mod go.sum embed.go ./
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/ cmd/sso/main.go

# Deploy the application binary into a lean image
FROM ubuntu:22.04 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/bin/main /app/main

EXPOSE ${REST_PORT}
EXPOSE ${GRPC_PORT}

ENTRYPOINT [ "./app/main" ] 