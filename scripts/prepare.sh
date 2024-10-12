#!/bin/bash

set -a
source .env
set +a

while [ $# -gt 0 ]; do
  case $1 in
        -s*|--service*)
            service="${1:2}"
            if [ -z "$service" ]; then
                echo "ERR: arg is empty"
                exit 1
            fi
            service_identified=true
        ;;
        *)
            echo "ERR: option $1 is not provided"
            exit 1
        ;;
    esac
    shift
done

if [[ $service_identified != true ]]; then
    echo "ERR: option --service required"
    exit 1
fi

echo creating dir for $service gen files
mkdir "$GOLANG_GEN_DIR/$service"

echo creating proto template for $service gen files
touch "$PROTOS_DIR/$service.proto"
