#!/bin/bash

set -e

# Build sso binary
go build \
    -o bin/ \
    -ldflags "\
    -X 'github.com/vishenosik/gocherry.BuildDate=$(date -u +'%Y-%m-%d-%H:%M:%S')' \
    -X 'github.com/vishenosik/gocherry.GitBranch=$(git branch --show-current)' \
    -X 'github.com/vishenosik/gocherry.GitCommit=$(git rev-parse HEAD)' \
    -X 'github.com/vishenosik/gocherry.GoVersion=$(go version)' \
    -X 'github.com/vishenosik/gocherry.GitTag=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "none")'"\
    \
    ./cmd/sso/main.go
