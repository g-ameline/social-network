#!/usr/bin/env bash

# echo " building data server"
# docker image build --progress=plain --file ../../containerizing/data_server/Dockerfile --tag sonet-database_server-image ../../
# echo " building app server"
# docker image build --progress=plain --file ../../containerizing/app_server/Dockerfile --tag sonet-application_server-image ../../


# echo " running containerized data server"
# docker container run --detach --publish 5000:5000 --name sonet-database_server-container sonet-database_server-image
# echo " running containerized app server"
# docker container run --detach --network=host --publish 8000:8000 --name sonet-application_server-container sonet-application_server-image


# since we have 3 server to run, to do it from one script we start it in background
# but then they do not stop with the script so we use the trap command
trap 'kill $(jobs -p)' EXIT
# the ampersand at the end makes it run in background ~= async
sh -c 'cd ./data_server/ && bash run_docker.sh' &
sh -c 'cd ./web_server/ && bash run_docker.sh' &
sh -c 'cd ./app_server/ && bash run_docker.sh' 

# trap 'trap - SIGTERM && kill 0' SIGINT SIGTERM EXIT
