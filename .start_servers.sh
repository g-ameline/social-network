#!/usr/bin/env bash
# since we have 3 server to run, to do it from one script we start it in background
# but then they do not stop with the script so we use the trap command
trap 'kill $(jobs -p)' EXIT
trap 'cd ./presentation/ && caddy stop' EXIT
# the ampersand at the end makes it run in background ~ parallel
sh -c 'cd ./presentation/ && caddy start' &
sh -c 'cd ./data/serve && go run . -- quiet' &
sh -c 'cd ./logic/ && go run . -- quiet' 

# trap 'trap - SIGTERM && kill 0' SIGINT SIGTERM EXIT
