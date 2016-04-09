#!/bin/bash
set -e

docker-compose kill
docker-compose rm -f
docker-compose up -d
cd www
gulp watch
#gulp_pid = $!

