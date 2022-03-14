#!/bin/bash

if docker network ls | grep "vscode-mr-plow"; then
  echo "network vscode-mr-plow already exists, nothing to do"
else
  echo "network vscode-mr-plow doesn't exist, creating ..."
  docker network create --driver=bridge --attachable --subnet=10.70.67.0/24 vscode-mr-plow
fi