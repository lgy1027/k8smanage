#!/bin/bash
ID=`ps -ef | grep kubemanage | grep -v "$0" | grep -v "grep" | awk '{print $2}'`
echo $ID
for id in $ID; do
    kill -9 $id
    echo "killed $id"
done

cd /opt/lgy/k8smanage
nohup /opt/lgy/k8smanage/kubemanage >> /tmp/k8smanage.log 2>&1 &