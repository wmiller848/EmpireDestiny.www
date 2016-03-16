#!/bin/bash
set -e

docker-compose build
mkdir -p ./bin
rm -rf ./bin/*
docker run -v `pwd`/bin:/cbin www_www_build cp -R /opt/empire.www/bin/. /cbin
