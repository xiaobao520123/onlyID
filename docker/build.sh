#!/bin/bash

go build -o onlyID ..

echo -n "Enter Docker hub username: "
read DOCKER_HUB_USERNAME
echo -n "Tag: "
read TAG
docker build -t $DOCKER_HUB_USERNAME/onlyid:$TAG .
docker push $DOCKER_HUB_USERNAME/onlyid:$TAG
