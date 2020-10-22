#!/bin/bash
set -e
EXE_FILE_NAME=kubemanage
HOST_LIST=(order)
HOST_PATH=/opt/lgy/k8smanage

######
# 编译
echo "1. 开始编译"

rm -rf bin
mkdir bin

GOOS=linux GOARCH=amd64 go build -o bin/${EXE_FILE_NAME}
echo -e "\
md5         $(md5sum bin/${EXE_FILE_NAME} | awk '{print $1}')\n\
git-hash    $(git log --pretty=format:"%H" -n 1)\n\
date        $(date "+%Y-%m-%d %H:%M:%S")\n\
hostname    $(hostname)\n\
git-user    $(git config --get user.name)\n\
" > bin/deploy.txt

######
# 打包
echo "2. 开始打包"

tar -zcf bin.tar.gz bin/

######
# 部署
echo "3. 开始部署"

for host in "${HOST_LIST[@]}"
do
    echo "  ==> ${host}"

    scp bin.tar.gz root@${host}:${HOST_PATH} 2>&1 >/dev/null
    scp -r ./config root@${host}:${HOST_PATH} 2>&1 >/dev/null
    scp -r ./docs root@${host}:${HOST_PATH} 2>&1 >/dev/null

    ssh root@${host} "\
    cd ${HOST_PATH}; \
    mv bin bin.$(date "+%Y%m%d%H%M%S"); \
    ps aux | grep -v grep | grep ${EXE_FILE_NAME} | awk '{print \$2}' | xargs -I {} kill {}; \
    sleep 3; \
    tar -zxf bin.tar.gz &> /dev/null; \
    cp ./bin/${EXE_FILE_NAME} .; \
    chmod a+x ${EXE_FILE_NAME}; \
    rm  bin.tar.gz; \
    nohup ./${EXE_FILE_NAME} >> log.$(date "+%Y%m%d%H%M%S").log &"
done

rm bin.tar.gz