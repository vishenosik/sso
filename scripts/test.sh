#!/bin/bash -e

# empty file
if [ -s ./protos/authorization.proto ]; then
        echo full
else
        echo empty
fi

# empty dir
if ! [[ -d ./protos ]]; then
        echo full
else
        echo empty
fi