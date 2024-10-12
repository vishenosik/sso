#!/bin/bash

set -a
source .env
set +a

for file in $PROTOS_DIR/*.proto; do

    service=$(basename "${file%.*}" .proto)
    dir="$GOLANG_GEN_DIR/"$service""

    echo "Generating $service service files"
    protoc -I $PROTOS_DIR "$PROTOS_DIR/$service.proto" --go_out=$dir --go_opt=paths=source_relative --go-grpc_out=$dir --go-grpc_opt=paths=source_relative

done

echo "Success"