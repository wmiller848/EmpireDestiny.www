#!/bin/bash
set -e

docker-compose rm -f
docker-compose up
mkdir -p `pwd`/bin
rm -rf `pwd`/bin/*
docker cp www_artifact_1:/opt/empire.www/bin/. `pwd`/bin/
docker exec empiredestiny_nginx_1 cp -R /opt/empire.www/. /usr/share/nginx/html/
