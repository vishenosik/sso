#!/bin/bash

set -a
source .env
set +a

service=$1

if [ -z "$service" ]; then
    echo "ERR: service name not provided"
    exit 1
fi

echo creating dir for $service gen files
mkdir "$GOLANG_GEN_DIR/$service"

echo creating proto template for $service gen files
touch "$PROTOS_DIR/$service.proto"

echo "Success"