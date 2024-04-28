#!/bin/bash -e
DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )
export PATH=$PATH:${DIR}/webui/bin
exec ${DIR}/webui/bin/yt-dlp-webui \
  -conf /var/snap/peerube/current/config/webui.yaml  \
  -db /var/snap/peerube/current/local.db

