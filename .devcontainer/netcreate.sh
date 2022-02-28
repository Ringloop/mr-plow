#!/bin/bash

if docker network ls | grep "mr-plow"; then
  echo "mr-plow already exist"
else
  echo "mr-plow network doesn't exist"
fi