#!/usr/bin/env bash
# remove contianers
docker rm sonet-web_server-container
docker rm sonet-application_server-container
docker rm sonet-database_server-container
# remove images
docker rmi Image sonet-web_server-image
docker rmi Image sonet-application_server-image
docker rmi Image sonet-database_server-image
# prune stuff
docker system prune
