#!/bin/bash
set -e

export GOPATH=${HOME}/.go
export GO111MODULE=auto
export GOPROXY=https://goproxy.cn

NAME="dmkit"
VERSION="$1"

CGO_ENABLED=0 go build -a -installsuffix cgo -o bin/${NAME} ./app.go
chmod +x bin/${NAME} 
