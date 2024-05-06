#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
export NODE_ENV=production
export NODE_CONFIG_DIR=/var/snap/peertube/current/config
export PATH=$PATH:${DIR}/peertube/bin
cd ${DIR}/peertube/app
/bin/rm -f /var/snap/peertube/current/peertube.socket
exec node dist/server
