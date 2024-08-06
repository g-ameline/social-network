#!/usr/bin/env bash
echo " building web server"
docker image build --progress=plain --file ../../containerizing/web_server/Dockerfile --tag sonet-web_server-image ../../
# docker image build --progress=plain --no-cache --file ../../containerizing/web_server/Dockerfile --tag sonet-web_server-image ../../
echo " running containerized web server"
# docker container run --detach --publish 3000:3000 --rm --name sonet-web_server-container sonet-web_server-image
docker container run --network=host --publish 3000:3000 --name sonet-web_server-container sonet-web_server-image
