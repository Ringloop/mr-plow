#!/bin/bash

WRK_PATH=$1
REL_PATH=.devcontainer

$WRK_PATH/$REL_PATH/gethnfo.sh $WRK_PATH
$WRK_PATH/$REL_PATH/netcreate.sh $WRK_PATH