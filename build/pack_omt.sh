#!/bin/sh
set -e

VERSION=$1
DATE=$(date +%Y%m%d%H)

SCHEMA=$2
TYPE=$3
TYPE_NAME=
TYPE_FILE=""
CGI_MAKETYPE=""

case $TYPE in
a2)
  TYPE_PLATFORM="a2"
  TYPE_NAME="AU_A302"
  TYPE_FILE="dasau_type.txt"
  CGI_MAKETYPE="zynq"
  ;;
a3)
  TYPE_PLATFORM="a3"
  TYPE_NAME="AU_A402"
  TYPE_FILE="dasa3_type.txt"
  CGI_MAKETYPE="zynqMp"
  ;;
e2)
  TYPE_PLATFORM="a2"
  TYPE_NAME="EU_E312"
  TYPE_FILE="das_e2.txt"
  CGI_MAKETYPE="zynq"
  ;;
e3)
  TYPE_PLATFORM="a3"
  TYPE_NAME="EU_E412"
  TYPE_FILE="das_e3.txt"
  CGI_MAKETYPE="am335x"
  ;;
eu)
  TYPE_PLATFORM="au"
  TYPE_NAME="EU_E212"
  TYPE_FILE="daseu_type.txt"
  CGI_MAKETYPE="am335x-oldsdk"
  ;;
eus)
  TYPE_PLATFORM="a2"
  TYPE_NAME="EU_E303"
  TYPE_FILE="daseus_type.txt"
  CGI_MAKETYPE="zynqMp"
  ;;
h2)
  TYPE_PLATFORM="a2"
  TYPE_NAME="H2RU_R311"
  TYPE_FILE="das_h2.txt"
  CGI_MAKETYPE="zynqMp"
  ;;
hp)
  TYPE_PLATFORM="au"
  TYPE_NAME="HRU_R211B"
  TYPE_FILE="das_hp.txt"
  CGI_MAKETYPE="am335x"
  ;;
m2)
  TYPE_PLATFORM="a2"
  TYPE_NAME="M2RU_R318"
  TYPE_FILE="das_m2.txt"
  CGI_MAKETYPE="zynqMp"
  ;;
m3h)
  TYPE_PLATFORM="a3"
  TYPE_NAME="M3RU-H_R416"
  TYPE_FILE="das_m3h.txt"
  CGI_MAKETYPE="zynqMp"
  ;;
m3l)
  TYPE_PLATFORM="a3"
  TYPE_NAME="M3RU-L_R417"
  TYPE_FILE="das_m3.txt"
  CGI_MAKETYPE="zynqMp"
  ;;
n2)
  TYPE_PLATFORM="a2"
  TYPE_NAME="N2RU_R304"
  TYPE_FILE="das_n2.txt"
  CGI_MAKETYPE="zynq"
  ;;
n3)
  TYPE_PLATFORM="a3"
  TYPE_NAME="N3RU_R404"
  TYPE_FILE="das_n3.txt"
  CGI_MAKETYPE="zynqMp"
  ;;
x2)
  TYPE_PLATFORM="a2"
  TYPE_NAME="X2RU_R328"
  TYPE_FILE="das_x2.txt"
  CGI_MAKETYPE="zynq"
  ;;
  *)
  echo "invalid type: ${TYPE}"
  exit -1
esac

PACK_DIR=$(cd "$(dirname "$0")"; pwd)
BIN_PATH="${PACK_DIR}/bin"
APP_DIST_PATH="${PACK_DIR}/dist"
INIT_PATH="${PACK_DIR}/init"
CGI_PATH="${PACK_DIR}/${SCHEMA}/${TYPE_PLATFORM}/cgi/${CGI_MAKETYPE}"

LOCAL_DIR=$(pwd)
TEMP_PATH="${LOCAL_DIR}/temp"
OUT_PATH="${LOCAL_DIR}/out"
NAME="${TYPE_NAME}_WEBOMT_${VERSION}_${DATE}"
OUTPUT_FILE="${OUT_PATH}/${NAME}.zip"

echo "name: $NAME"
echo "version: $VERSION"

mkdir -p ${OUT_PATH}
rm -rf ${TEMP_PATH}
mkdir -p ${TEMP_PATH}

cp ${INIT_PATH}/update.sh ${TEMP_PATH}/update.sh
cp ${INIT_PATH}/${TYPE_FILE} ${TEMP_PATH}/type.txt

mkdir -p ${TEMP_PATH}/webomt/bin
cp -r ${BIN_PATH}/gomt.arm ${TEMP_PATH}/webomt/bin/gomt

mkdir -p ${TEMP_PATH}/webomt/web/www
cp -r ${APP_DIST_PATH}/* ${TEMP_PATH}/webomt/web/www/
cp -r ${CGI_PATH}/cgi-bin ${TEMP_PATH}/webomt/web/www/cgi-bin

cd ${TEMP_PATH}/webomt/web
find www -type f -print0 | xargs -0 md5sum | sort > websum.md5

chmod a+x ${TEMP_PATH}/update.sh
chmod a+r ${TEMP_PATH}/type.txt
chmod -R a+r ${TEMP_PATH}/webomt
chmod -R a+x ${TEMP_PATH}/webomt/web/www/cgi-bin

cd ${TEMP_PATH}
zip -r -P SunWave321 ${OUTPUT_FILE} ./*
cd ../
rm -rf ${TEMP_PATH}

echo "package: ${OUTPUT_FILE}"
echo ""
