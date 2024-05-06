#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
export NODE_ENV=production
export NODE_CONFIG_DIR=/var/snap/peertube/current/config
export PATH=$PATH:${DIR}/peertube/bin
cd ${DIR}/peertube/app
exec ${DIR}/peertube/bin/node.sh dist/server

