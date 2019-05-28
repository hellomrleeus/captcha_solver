#!/bin/bash
COUNT=0
docker rm 1mclient_$COUNT
while [ $? -eq 0 ]; do
    COUNT=$(($COUNT+1))
    docker rm 1mclient_$COUNT
done
