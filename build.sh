#!/bin/bash
export GO111MODULE=on
projectpath=`pwd`
#export GOPATH="${projectpath}:$GOPATH"
export GOBIN="${projectpath}/bin"
echo GOPATH=$GOPATH
echo pwd=$PWD
go build -o engine plugin.test/main/cmd/engine

if test $? -gt 0;then echo "Error to build";exit 12;fi


enginePs="`ps -f|grep engine|grep -v grep`"
if ! test -z "$enginePs" ;then 
	pid=$(echo $enginePs|cut -d ' ' -f 2) #前面还有一个UID,所以
	echo "enginePS=$enginePs pid=$pid"
	kill -9 $pid
fi
rm -f nohup.out
nohup ./engine& 2>/dev/null

sleep 1
echo ">> show the plugin value..."
curl "http://localhost:12345/hello"

tail -f nohup.out

