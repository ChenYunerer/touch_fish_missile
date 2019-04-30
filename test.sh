#!/usr/bin/env bash
#测试脚本
echo "编译开始"
filename="main"
go build -o $filename ./src/main/
echo "编译完成"
echo "启动Server"
./$filename -startType server > ./output/server.log 2>&1 &
echo "输入Client数量"
read numOfClient
num=0
while ((num < numOfClient))
do
    echo "启动Client$num"
    ./$filename -startType client -token 10${num} > ./output/client10${num}.log 2>&1 &
    let "num++"
done
echo "启动完毕 可通过tail -f output路径下的文件查看日志"
echo "通过killall ${filename}来关闭所有相关进程"


