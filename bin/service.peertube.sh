#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
export PATH=$PATH:${DIR}/webui/bin
cd ${DIR}/peertube/app
exec ${DIR}/peertube/bin/node.sh dist/server

