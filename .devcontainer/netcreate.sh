#!/bin/bash

if docker network ls | grep "mr-plow"; then
  echo "mr-plow network already exist"
else
  echo "mr-plow network doesn't exist, creating.."
  docker network create --driver=bridge --attachable --subnet=10.70.67.0/24 mr-plow
fi