#!/bin/bash


#!/bin/bash
ENV_FILE=".env"
CMD=${@:2}

set -a # automatically export all variables
source .env
set +a

echo $CMD

# source '.env' && bash 'scriptname.sh'


# echo $PROTOS_DIR

# mkdir ./gen/go/authentication
# echo "DONE"