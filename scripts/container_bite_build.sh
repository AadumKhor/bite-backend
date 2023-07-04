#!/usr/bin/env bash
set -o errexit

sudo docker build -t bite-server:$(git rev-parse --short HEAD) -f ./dockerfiles/Dockerfile .