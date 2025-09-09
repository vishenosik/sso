#!/bin/bash

OUTPUT_DIR="dist"

PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)


mkdir -p $OUTPUT_DIR

for platform in "${PLATFORMS[@]}"; do

    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    output_name=$OUTPUT_DIR/${GOOS}-${GOARCH}/sso
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building for $GOOS/$GOARCH..."
  
    # Build sso binary
    env GOOS=$GOOS GOARCH=$GOARCH go build \
        -o $output_name \
        -ldflags "\
        -X 'main.BuildDate=$(date -u +'%Y-%m-%d-%H:%M:%S')' \
        -X 'main.GitBranch=$(git branch --show-current)' \
        -X 'main.GitCommit=$(git rev-parse HEAD)' \
        -X 'main.GoVersion=$(go version)' \
        -X 'main.GitTag=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "none")'"\
        \
        ./cmd/sso/main.go

done