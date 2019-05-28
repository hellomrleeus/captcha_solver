#!/bin/bash
CONNECTIONS=$1
REPLICAS=$2
IP=$3
#go build --tags "static netgo" -o client client.go
for (( c=0; c<${REPLICAS}; c++ ))
do
    docker run -v $(pwd)/client:/client --name 1mclient_$c -d alpine /ls -la client \
    # -conn=${CONNECTIONS} -ip=${IP}
done