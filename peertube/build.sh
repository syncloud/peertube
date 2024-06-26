#!/bin/sh -xe

DIR=$( cd "$( dirname "$0" )" && pwd )
cd ${DIR}

apt update
apt install -y wget

BUILD_DIR=${DIR}/../build/snap/peertube
mkdir -p ${BUILD_DIR}
sed -i 's/server\.listen(port, hostname/server\.listen\(hostname/g' /app/dist/server.js
sed -i 's#new transports.File.*#// no file#g' /app/dist/core/helpers/logger.js
sed -i 's#new transports.File#new transports.Console#g' /app/dist/core/helpers/audit-logger.js
mv /app ${BUILD_DIR}/app
cp -r /opt ${BUILD_DIR}
cp -r /usr ${BUILD_DIR}
cp -r /lib ${BUILD_DIR}

ldd ${BUILD_DIR}/usr/bin/ffmpeg
mkdir ${BUILD_DIR}/bin
cp $DIR/bin/* ${BUILD_DIR}/bin
ldd ${BUILD_DIR}/usr/bin/ffmpeg

wget https://framagit.org/framasoft/peertube/official-plugins/-/archive/master/official-plugins-master.tar.gz
tar xf official-plugins-master.tar.gz
mkdir ${BUILD_DIR}/app/plugins
sed -i 's/role = parseInt\(.*\)/role = roleToParse == "syncloud" ? 0: 2/' official-plugins-master/peertube-plugin-auth-openid-connect/main.js
cp -r official-plugins-master/peertube-plugin-auth-openid-connect ${BUILD_DIR}/app/plugins

