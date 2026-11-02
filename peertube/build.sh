#!/bin/sh -xe

DIR=$( cd "$( dirname "$0" )" && pwd )
cd ${DIR}

OIDC_PLUGIN=2.1.1
YT_DLP=2026.10.26

${DIR}/../apt.sh wget

BUILD_DIR=${DIR}/../build/snap/peertube
mkdir -p ${BUILD_DIR}
sed -i 's/server\.listen(port, hostname/server\.listen\(hostname/g' /app/dist/server.js
grep -q 'server.listen(hostname' /app/dist/server.js
sed -i '/new transports.File(fileLoggerOptions)/d; /loggerTransports.push(fileTransport)/d' /app/dist/core/helpers/logger.js
! grep -q 'transports.File' /app/dist/core/helpers/logger.js
sed -i 's#new transports.File#new transports.Console#g' /app/dist/core/helpers/audit-logger.js
rm -f /usr/local/lib/node_modules/pnpm/pn /usr/local/lib/node_modules/pnpm/pnpx /usr/local/lib/node_modules/pnpm/pnx
rm -f /usr/local/bin/pn /usr/local/bin/pnpx /usr/local/bin/pnx
mv /app ${BUILD_DIR}/app
cp -r /opt ${BUILD_DIR}
cp -r /usr ${BUILD_DIR}
cp -r /lib ${BUILD_DIR}

ldd ${BUILD_DIR}/usr/bin/ffmpeg
mkdir ${BUILD_DIR}/bin
cp $DIR/bin/* ${BUILD_DIR}/bin
ldd ${BUILD_DIR}/usr/bin/ffmpeg

wget --tries=3 https://registry.npmjs.org/peertube-plugin-auth-openid-connect/-/peertube-plugin-auth-openid-connect-${OIDC_PLUGIN}.tgz
tar xzf peertube-plugin-auth-openid-connect-${OIDC_PLUGIN}.tgz
sed -i 's/role = parseInt\(.*\)/role = roleToParse == "syncloud" ? 0: 2/' package/dist/main.js
grep -q 'role = roleToParse == "syncloud"' package/dist/main.js
cd package
npm install --omit=dev --no-audit --no-fund
node -e "require('./dist/main.js')"
cd ${DIR}
mkdir ${BUILD_DIR}/app/plugins
cp -r package ${BUILD_DIR}/app/plugins/peertube-plugin-auth-openid-connect

YT_DLP_BINARY=yt-dlp_linux
if [ "$(dpkg --print-architecture)" = "arm64" ]; then
  YT_DLP_BINARY=yt-dlp_linux_aarch64
fi
mkdir ${BUILD_DIR}/yt-dlp
wget --tries=3 https://github.com/yt-dlp/yt-dlp/releases/download/${YT_DLP}/${YT_DLP_BINARY} -O ${BUILD_DIR}/yt-dlp/yt-dlp
chmod +x ${BUILD_DIR}/yt-dlp/yt-dlp
