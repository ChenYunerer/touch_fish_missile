#!/usr/bin/env bash
#linux平台构建脚本 ./build.sh main darwin/linux
#默认目标名称为main
name=$1
os=$2
if [ ${#name} = 0 ]
then
name="main"
fi
if [ ${#os} = 0 ]
then
os="darwin"
fi
CGO_ENABLED=0 GOOS=${os} GOARCH=amd64 go build -o ${name} ./src/main/