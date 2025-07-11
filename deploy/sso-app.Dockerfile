# syntax=docker/dockerfile:1

FROM golang:1.24.2 AS build-stage

# Set destination for COPY
WORKDIR /app

COPY . .

RUN go mod download

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o bin/ cmd/sso/main.go

# Deploy the application binary into a lean image
FROM ubuntu:22.04 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/bin/main /app/main

# HTTP port
EXPOSE 4080
# gRPC port
EXPOSE 4090

ENTRYPOINT [ "./app/main" ] 