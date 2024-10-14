#!/bin/bash -e
if [ -s ./protos/authorization.proto ]; then
        echo full
else
        echo empty
fi