#!/bin/bash
set -e

cd deployments

docker-compose -p db-tests -f docker-compose.dev.yaml build --no-cache db-test
docker-compose -p db-tests -f docker-compose.dev.yaml run db-test
docker-compose -p db-tests -f docker-compose.dev.yaml down