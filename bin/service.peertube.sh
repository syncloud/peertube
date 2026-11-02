#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
export NODE_ENV=production
export NODE_CONFIG_DIR=/var/snap/peertube/current/config
export NODE_EXTRA_CA_CERTS=/var/snap/platform/current/syncloud.ca.crt
export PATH=$PATH:${DIR}/peertube/bin
cd ${DIR}/peertube/app
/bin/rm -f /var/snap/peertube/current/peertube.socket
$DIR/bin/wait-for-configure.sh
exec node dist/server
