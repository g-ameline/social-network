#!/usr/bin/env bash
echo " building app server"
docker image build --progress=plain --file ../../containerizing/app_server/Dockerfile --tag sonet-application_server-image ../../
# docker image build --progress=plain --no-cache --file ../../back_end/containerizing/Dockerfile --tag sonet-application_server-image ../../
echo " running containerized app server"
# docker container run --detach --publish 8000:8000 --rm --name sonet-application_server-container sonet-application_server-image
docker container run --network=host --publish 8000:8000 --name sonet-application_server-container sonet-application_server-image
