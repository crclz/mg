#!/bin/bash

set -e
# set -x # enable this to print details

ServicePsm="p.s.m"
ServiceCluster="default"
RuntimeIdcName="hl"

MeshLoaderImage="hub.byted.org/mesh/mesh_loader:20231027"
GoTestOutputFile="gotest.out"

if [[ "$CompilePackage" == "" ]];then
   echo "env var is empty: CompilePackage"
   exit 1
fi

if [[ "$TestRunPattern" == "" ]];then
   echo "env var is empty: TestRunPattern"
   exit 1
fi

if [[ "$DockerDisableInteractive" == "1" ]]; then
   interactiveOption=""
else
   interactiveOption="-t"
fi


# build
go test -c -o $GoTestOutputFile --gcflags "all=-l -N" $CompilePackage

# run test
doas -p $ServicePsm docker run -i --rm \
   $interactiveOption\
   --network=host -e CONSUL_HTTP_HOST=`hostname -i` \
   -e RUNTIME_IDC_NAME=$RuntimeIdcName  -e RUNTIME_SERVICE_PORT=1234  -e RUNTIME_DEBUG_PORT=11234 \
   -e PSM=$ServicePsm -e SERVICE_CLUSTER=$ServiceCluster \
   -e `doas -p $ServicePsm env | grep SEC_TOKEN` \
   -v $(pwd):/workdir \
   $MeshLoaderImage \
   mesh_loader --psm $ServicePsm \
   --cluster $ServiceCluster \
   --ipv6 0:0:0:0::0 \
   --erpc \
   -c "./$GoTestOutputFile --test.v --test.run $TestRunPattern" \

# 20231027版本的镜像 需要添加 --ipv6 xxx 使用