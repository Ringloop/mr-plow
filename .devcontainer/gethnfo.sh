#!/bin/bash

WRK_PATH=$1
REL_PATH=.devcontainer/dockerfiles/mrplow/hnfo

id -u > $WRK_PATH/$REL_PATH
id -g >> $WRK_PATH/$REL_PATH
getent group docker | cut -d: -f3 >> $WRK_PATH/$REL_PATH
echo "provisioning with host's user UID: $(id -u), GID: $(id -g), host's docker GID: $(getent group docker | cut -d: -f3)"