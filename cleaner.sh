#!/bin/bash
docker image prune -a -f
docker container prune -f
docker builder prune -f
docker system df
