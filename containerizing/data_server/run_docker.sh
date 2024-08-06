#!/usr/bin/env bash
echo " building data server"
docker image build --progress=plain --file ../../containerizing/data_server/Dockerfile --tag sonet-database_server-image ../../
# docker image build --progress=plain --no-cache --file ../../back_end/containerizing/Dockerfile --tag sonet-database_server-image ../../
echo " running containerized data server"
# docker container run --detach --publish 5000:5000 --rm --name sonet-database_server-container sonet-database_server-image
docker container run  --network=host --publish 5000:5000 --name sonet-database_server-container sonet-database_server-image
