#!/bin/bash

GOLANG_GEN_DIR=$1
PROTOS_DIR=$2

for VersionDir in $PROTOS_DIR/*; do

    ApiVersion=$(basename "${VersionDir}")

    for file in $VersionDir/*.proto; do

        service=$(basename "${file%.*}" .proto)
        dir="$GOLANG_GEN_DIR/$ApiVersion"

        if ! [[ -d $dir ]]; then
            mkdir $dir
        fi

        dir="$GOLANG_GEN_DIR/$ApiVersion/$service"
        if ! [[ -d $dir ]]; then
            mkdir $dir
        fi

        echo "Generating service files: $service"
        protoc -I $PROTOS_DIR/$ApiVersion "$PROTOS_DIR/$ApiVersion/$service.proto" --go_out=$dir --go_opt=paths=source_relative --go-grpc_out=$dir --go-grpc_opt=paths=source_relative

    done

done

echo "Success"