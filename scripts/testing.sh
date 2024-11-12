#!/bin/bash

COVERAGE_FILE=$1
TESTING_LIST=$(go list ./... | grep -E "internal|pkg" | grep -v mocks )

go test -v -cover -coverprofile=$COVERAGE_FILE $TESTING_LIST