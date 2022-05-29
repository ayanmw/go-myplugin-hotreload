#!/bin/bash
export GO111MODULE=on
projectpath=`pwd`
#export GOPATH="${projectpath}:$GOPATH"
export GOBIN="${projectpath}/bin"
pluginname="plugin$1"
pluginpath="plugin.test/main/cmd/logicv$1"

GOVAR1="-X 'plugin.test/main/cmd/logicv$1.ModuleVersion=100$1' "
GOVAR2="-X 'plugin.test/main/cmd/logicv$1.ModuleName=game_`date '+%Y-%m-%d %H:%M:%S'`_$1' "
#    -ldflags="-X 'plugin.test/main/cmd/logicv$1.ModuleVersion=100$1' -X 'plugin.test/main/cmd/logicv$1.ModuleName=game_`date '+%Y-%m-%d %H:%M:%S'`_$1'" \

cd cmd
mv logic logicv$1
cd -
#    -ldflags="-pluginpath=$pluginpath" \ ##如果跟so的包名不一样, 反倒是 找不到任何 symbol,所以自定义pluginpath 反而没有意义了.
go build -buildmode plugin -o $pluginname.so \
    -ldflags="$GOVAR1 $GOVAR2"\
    plugin.test/main/cmd/logicv$1

cd cmd
mv logicv$1 logic
cd -


echo ">> hot reload plugin..."
curl "http://localhost:12345/load?name=$pluginname.so"
curl "http://localhost:12345/gc"
echo ">> show the plugin value..."
curl "http://localhost:12345/hello"

