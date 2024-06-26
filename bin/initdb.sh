#!/bin/bash -e

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )

if [ -z "$SNAP_COMMON" ]; then
  echo "SNAP_COMMON environment variable must be set"
  exit 1
fi

# shellcheck source=config/env
. "${SNAP_DATA}/config/env"

if [[ "$(whoami)" == "peertube" ]]; then
    ${DIR}/postgresql/bin/initdb.sh "$@"
else
    sudo -E -H -u peertube ${DIR}/postgresql/bin/initdb.sh "$@"
fi
