#!/bin/bash

ARGS=$1

export REPOPATH=$PWD

docker build -f $REPOPATH/docker/go/Dockerfile -t go .
docker run -i -v $PWD:/root/rccmd go /root/rccmd/compile.sh

if [ $? -ne 0 ]; then
	exit 1
fi

#docker stack deploy --compose-file $REPOPATH/docker/docker-compose.yml tcp

if [ $ARGS == "debug" ]; then
	docker compose -f docker/docker-compose.yml up
else
	docker compose -f docker/docker-compose.yml up --exit-code-from sut
fi
