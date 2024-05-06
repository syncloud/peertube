#!/bin/sh -ex

DIR=$( cd "$( dirname "$0" )" && pwd )
cd ${DIR}

BUILD_DIR=${DIR}/../build/snap/peertube
${BUILD_DIR}/bin/node.sh --version

$BUILD_DIR/bin/ffmpeg --help
$BUILD_DIR/bin/ffprobe --help
