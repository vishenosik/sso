#!/bin/bash

set -e

# Build sso binary
go build -o bin/ -ldflags "\
    -X 'main.BuildDate=$(date -u +'%Y-%m-%d-%H:%M:%S')' \
    -X 'main.GitBranch=$(git branch --show-current)' \
    -X 'main.GitCommit=$(git rev-parse HEAD)' \
    -X 'main.GoVersion=$(go version)' \
    -X 'main.GitTag=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "none")'"\
    \
    ./cmd/sso/main.go
