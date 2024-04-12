#!/bin/sh -e

DIR=$( cd "$( dirname "$0" )" && pwd )
cd ${DIR}

VERSION=$1
BUILD_DIR=${DIR}/../build/snap/webui
mkdir -p $BUILD_DIR/bin

cd ${DIR}/../build
wget https://github.com/Chocobozzz/PeerTube/releases/download/v$VERSION/peertube-v$VERSION.tar.xz
tar xf peertube-v$VERSION.tar.xz
cp -r ${DIR}/bin/* ${BUILD_DIR}/bin