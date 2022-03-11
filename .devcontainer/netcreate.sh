#!/bin/bash

if docker network ls | grep "mr-plow"; then
  echo "network mr-plow already exists, nothing to do"
else
  echo "network mr-plow doesn't exist, creating ..."
  docker network create --driver=bridge --attachable --subnet=10.70.67.0/24 mr-plow
fi