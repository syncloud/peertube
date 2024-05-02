#!/bin/sh -xe

DIR=$( cd "$( dirname "$0" )" && pwd )
cd ${DIR}

BUILD_DIR=${DIR}/../build/snap/peertube
mkdir -p ${BUILD_DIR}
sed -i 's/server\.listen(port, hostname/server\.listen\(hostname/g' /app/dist/server.js
sed -i 's#new transports.File.*#// no file#g' /app/dist/core/helpers/logger.js
mv /app ${BUILD_DIR}/app
cp -r /opt ${BUILD_DIR}
cp -r /usr ${BUILD_DIR}
cp -r /bin ${BUILD_DIR}
cp -r /lib ${BUILD_DIR}
cp $DIR/bin/* ${BUILD_DIR}/bin
