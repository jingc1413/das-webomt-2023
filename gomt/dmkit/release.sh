set -e
rm -rf out

OUT_DIR=../../models

bin/dmkit generate --import ../../asserts/device-files
rm -rf ${OUT_DIR} 
cp -r out ${OUT_DIR} 

find ${OUT_DIR}/mock -name "*.csv" | sed 's/.*/"&"/' | xargs rm -f
find ${OUT_DIR}/mock -name "*.gz" | sed 's/.*/"&"/' | xargs rm -f

APP_DIR=../../app-v3
mkdir -p ${APP_DIR}/public/mock
rm -rf ${APP_DIR}/public/mock/models
cp -r ${OUT_DIR}/mock/corning ${APP_DIR}/public/mock/models


# bin/dmkit generate --import ../../asserts/device-files/a3-corning-20240509
# rm -rf ~/coding/das-webomt/app-v3/public/models
# cp -r out ~/coding/das-webomt/app-v3/public/models

# find ~/coding/das-webomt/app-v3/public/models -name "*.js" | sed 's/.*/"&"/' | xargs rm -f
# find ~/coding/das-webomt/app-v3/public/models -name "*.csv" | sed 's/.*/"&"/' | xargs rm -f
# find ~/coding/das-webomt/app-v3/public/models -name "*.json" | sed 's/.*/"&"/' | xargs rm -f
