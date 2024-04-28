#!/bin/sh -xe

DIR=$( cd "$( dirname "$0" )" && pwd )
cd ${DIR}

VERSION=$1
BUILD_DIR=${DIR}/../build/snap/peertube
mkdir -p ${BUILD_DIR}
apt update
apt install -y wget xz-utils

cd ${DIR}/../build
wget https://github.com/Chocobozzz/PeerTube/releases/download/v$VERSION/peertube-v$VERSION.tar.xz
tar xf peertube-v$VERSION.tar.xz
rm -rf peertube-v$VERSION.tar.xz
cd peertube-v$VERSION
yarn install --production --pure-lockfile
cd ..
mv peertube-v$VERSION ${BUILD_DIR}/peertube
cp -r /opt ${BUILD_DIR}
cp -r /usr ${BUILD_DIR}
cp -r /bin ${BUILD_DIR}
cp -r /lib ${BUILD_DIR}
cp $DIR/bin/* ${BUILD_DIR}/bin
