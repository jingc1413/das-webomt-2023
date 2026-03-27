# clone stage
FROM dasomt/cgi:240620 AS clone-cgi
WORKDIR /source

# build stage
FROM linux_20231012:latest AS build-a3-cf-cgi
WORKDIR /source
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/gnu/arm-xilinx-linux-gnueabi-4.9.2/bin/
COPY --from=clone-cgi /source/a3-cf/cgi ./cgi
RUN cd /source/cgi && rm -rf cgi-bin && mkdir -p cgi-bin && make -f Makefile-zynq
RUN mkdir -p /source/out/zynq && mv /source/cgi/cgi-bin /source/out/zynq/cgi-bin
RUN cd /source/cgi && rm -rf cgi-bin && mkdir -p cgi-bin && make -f Makefile-zynqMp
RUN mkdir -p /source/out/zynqMp && mv /source/cgi/cgi-bin /source/out/zynqMp/cgi-bin
RUN cd /source/cgi && rm -rf cgi-bin && mkdir -p cgi-bin && make -f Makefile-am335x
RUN mkdir -p /source/out/am335x && mv /source/cgi/cgi-bin /source/out/am335x/cgi-bin
# RUN cd /source/cgi && rm -rf cgi-bin && mkdir -p cgi-bin && make -f Makefile-am335x-oldsdk
# RUN mkdir -p /source/out/am335x-oldsdk && mv /source/cgi/cgi-bin /source/out/am335x-oldsdk/cgi-bin

FROM linux_20231012:latest AS build-a3-cn-cgi
WORKDIR /source
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/gnu/arm-xilinx-linux-gnueabi-4.9.2/bin/
RUN rm -rf /source/cgi
COPY --from=clone-cgi /source/a3-cn/cgi /source/cgi
RUN cd /source/cgi && rm -rf cgi-bin && mkdir -p cgi-bin && make -f Makefile-zynq
RUN mkdir -p /source/out/zynq && mv /source/cgi/cgi-bin /source/out/zynq/cgi-bin
RUN rm -rf /source/cgi
COPY --from=clone-cgi /source/a3-cn/cgi /source/cgi
RUN cd /source/cgi && rm -rf cgi-bin && mkdir -p cgi-bin && make -f Makefile-zynqMp
RUN mkdir -p /source/out/zynqMp && mv /source/cgi/cgi-bin /source/out/zynqMp/cgi-bin
RUN rm -rf /source/cgi
COPY --from=clone-cgi /source/a3-cn/cgi /source/cgi
RUN cd /source/cgi && rm -rf cgi-bin && mkdir -p cgi-bin && make -f Makefile-am335x
RUN mkdir -p /source/out/am335x && mv /source/cgi/cgi-bin /source/out/am335x/cgi-bin
# RUN rm -rf /source/cgi
# COPY --from=clone-cgi /source/a3-cn/cgi /source/cgi
# RUN cd /source/cgi && rm -rf cgi-bin && mkdir -p cgi-bin && make -f Makefile-am335x-oldsdk
# RUN mkdir -p /source/out/am335x-oldsdk && mv /source/cgi/cgi-bin /source/out/am335x-oldsdk/cgi-bin

FROM node:20 AS build-node
RUN npm config set registry https://registry.npmmirror.com 
WORKDIR /source
COPY ./app-v3/package.json .
COPY ./app-v3/package-lock.json .
RUN npm install --force
COPY ./app-v3 .
RUN rm -rf ./public/mock
RUN rm -rf ./public/models
RUN rm -rf ./public/config
COPY ./models/models ./public/models
RUN find ./public/models -name "*.json" | sed 's/.*/"&"/' | xargs rm -f
RUN find ./public/models -name "*.csv" | sed 's/.*/"&"/' | xargs rm -f
RUN npm run build
RUN chmod -R a+r /source/dist

FROM golang:1.21 AS build-gomt
ENV GO111MODULE=auto
ENV GOPROXY=https://goproxy.cn
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=6
ENV CGO_ENABLED=0
ENV CC=arm-linux-gnueabi-gcc
ENV CXX=arm-linux-gnueabi-g++
ENV APP_NAME="gomt"
WORKDIR /source
COPY ./gomt .
RUN rm -rf go.work; go work init dmkit core gomt das
RUN cd gomt; rm -rf ./dist; mkdir -p dist; touch ./dist/empty
RUN cd gomt; go build -a -installsuffix cgo -ldflags "-s -w" -o /source/bin/${APP_NAME}.arm ./app.go

FROM golang:1.21 AS build-omtapp
ENV GO111MODULE=auto
ENV GOPROXY=https://goproxy.cn
ENV APP_NAME=omtapp
WORKDIR /source
COPY ./gomt .
RUN rm -rf ./gomt/dist
COPY --from=build-node /source/dist ./gomt/dist
RUN rm -rf go.work; go work init dmkit core gomt das
RUN cd gomt; GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o /source/bin/${APP_NAME}.exe ./app.go
RUN cd gomt; GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -o /source/bin/${APP_NAME} ./app.go

FROM base/alpine:3.18
ARG release_version=develop
WORKDIR /release
COPY --from=build-a3-cf-cgi /source/out/zynq  ./cf/a3/cgi/zynq
COPY --from=build-a3-cf-cgi /source/out/zynqMp  ./cf/a3/cgi/zynqMp
COPY --from=build-a3-cf-cgi /source/out/am335x  ./cf/a3/cgi/am335x
# COPY --from=build-a3-cf-cgi /source/out/am335x-oldsdk  ./cf/a3/cgi/am335x-oldsdk

COPY --from=build-a3-cn-cgi /source/out/zynq  ./cn/a3/cgi/zynq
COPY --from=build-a3-cn-cgi /source/out/zynqMp  ./cn/a3/cgi/zynqMp
COPY --from=build-a3-cn-cgi /source/out/am335x  ./cn/a3/cgi/am335x

COPY --from=build-node /source/dist ./dist
COPY --from=build-gomt /source/bin ./bin
COPY --from=build-omtapp /source/bin ./bin
COPY --from=clone-cgi /source/a3-cf/cgi/init ./init
COPY ./build/cgi_init_update.sh ./init/update.sh
RUN chmod a+x ./init/update.sh
COPY ./build/pack_omt.sh .
RUN mkdir -p pkgs
# RUN mv ./bin/gomt.arm ./bin/gomt.arm.orig && upx --brute -o ./bin/gomt.arm ./bin/gomt.arm.orig 

RUN rm -rf ./dist/public/config && mkdir -p ./dist/public/config
COPY ./asserts/device-files/a3-20231017/combo.ini ./dist/public/config/
COPY ./asserts/device-files/a3-20231017/xml/* ./dist/public/config/
RUN rm -rf ./dist/public/imgs && mkdir -p ./dist/public/imgs
RUN chmod -R a+r ./dist
# RUN ./pack_omt.sh ${release_version} eu
# RUN ./pack_omt.sh ${release_version} hp
# RUN ./pack_omt.sh ${release_version} a2
# RUN ./pack_omt.sh ${release_version} e2
# RUN ./pack_omt.sh ${release_version} h2
# RUN ./pack_omt.sh ${release_version} m2
# RUN ./pack_omt.sh ${release_version} n2
# RUN ./pack_omt.sh ${release_version} x2
RUN ./pack_omt.sh ${release_version} cf a3
RUN ./pack_omt.sh ${release_version} cf e3
# RUN ./pack_omt.sh ${release_version} m3h
# RUN ./pack_omt.sh ${release_version} m3l
RUN ./pack_omt.sh ${release_version} cf n3
# RUN ./pack_omt.sh ${release_version} eus
RUN mv out ./pkgs/sunwave

RUN rm -rf ./dist/public/config && mkdir -p ./dist/public/config
COPY ./asserts/device-files/a3-corning-20240326/combo.ini ./dist/public/config/
COPY ./asserts/device-files/a3-corning-20240326/xml/* ./dist/public/config/
RUN rm -rf ./dist/public/imgs && mkdir -p ./dist/public/imgs
COPY ./asserts/images/login-bg.corning.jpg ./dist/public/imgs/login-bg.jpg
RUN chmod -R a+r ./dist
# RUN ./pack_omt.sh ${release_version} eu
# RUN ./pack_omt.sh ${release_version} hp
# RUN ./pack_omt.sh ${release_version} a2
# RUN ./pack_omt.sh ${release_version} e2
# RUN ./pack_omt.sh ${release_version} h2
# RUN ./pack_omt.sh ${release_version} m2
# RUN ./pack_omt.sh ${release_version} n2
# RUN ./pack_omt.sh ${release_version} x2
RUN ./pack_omt.sh ${release_version} cn a3
RUN ./pack_omt.sh ${release_version} cn e3
# RUN ./pack_omt.sh ${release_version} m3h
# RUN ./pack_omt.sh ${release_version} m3l
RUN ./pack_omt.sh ${release_version} cn n3
# RUN ./pack_omt.sh ${release_version} eus
RUN mv out ./pkgs/corning

RUN rm -rf ./dist/public/config && mkdir -p ./dist/public/config
COPY ./asserts/device-files/a3-corning-demo/combo.ini ./dist/public/config/
COPY ./asserts/device-files/a3-corning-demo/xml/* ./dist/public/config/
RUN rm -rf ./dist/public/imgs && mkdir -p ./dist/public/imgs
COPY ./asserts/images/login-bg.corning.jpg ./dist/public/imgs/login-bg.jpg
RUN chmod -R a+r ./dist
# RUN ./pack_omt.sh ${release_version} eu
# RUN ./pack_omt.sh ${release_version} hp
# RUN ./pack_omt.sh ${release_version} a2
# RUN ./pack_omt.sh ${release_version} e2
# RUN ./pack_omt.sh ${release_version} h2
# RUN ./pack_omt.sh ${release_version} m2
# RUN ./pack_omt.sh ${release_version} n2
# RUN ./pack_omt.sh ${release_version} x2
RUN ./pack_omt.sh ${release_version} cn a3
RUN ./pack_omt.sh ${release_version} cn e3
# RUN ./pack_omt.sh ${release_version} m3h
# RUN ./pack_omt.sh ${release_version} m3l
RUN ./pack_omt.sh ${release_version} cn n3
# RUN ./pack_omt.sh ${release_version} eus
RUN mv out ./pkgs/corning-powersharing