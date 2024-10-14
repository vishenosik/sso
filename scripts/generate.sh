#!/bin/bash

set -a
source .env
set +a

for file in $PROTOS_DIR/*.proto; do

    service=$(basename "${file%.*}" .proto)
    dir="$GOLANG_GEN_DIR/$service"

    if ! [[ -d $dir ]]; then
        echo "Creating path: $dir"
        mkdir $dir
    fi

    echo "Generating service files: $service"
    protoc -I $PROTOS_DIR "$PROTOS_DIR/$service.proto" --go_out=$dir --go_opt=paths=source_relative --go-grpc_out=$dir --go-grpc_opt=paths=source_relative

done

echo "Success"