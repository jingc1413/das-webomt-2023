#!/bin/bash
set -e

export GOPATH=${HOME}/.go
export GO111MODULE=auto
export GOPROXY=https://goproxy.cn

export GOOS=linux
export GOARCH=arm
export GOARM=6
export CGO_ENABLED=0
export CC=arm-linux-gnueabi-gcc
export CXX=arm-linux-gnueabi-g++

NAME="gomt"

rm -f ./bin/${NAME}

echo "go build"
go build -a -installsuffix cgo -ldflags "-s -w" -o ./bin/${NAME}.arm ./app.go
# echo "upx"
# upx --brute -o ./bin/${NAME}.arm ./bin/${ORIG}

