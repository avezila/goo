#!/bin/bash
set -e
set -v
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/../..
cd "$ROOT"

if [[ $# < 2 ]]; then
	echo "usage:  build.sh service1 service2 ... image-name"
	exit 1
fi

imagename=${@: -1}
output=`echo $imagename | tr "/" __`
output=`echo $output | tr ":" _`
services=${@: 1 : $#-1}
allservices=$( ls -1 services)

protoReg='^(.*)\.proto$'


rm -rf `find . | grep "\.pb.js$"`

for service in $services; do
  for other in $allservices; do
    cd "$ROOT"/services/$other
    for protof in $( ls -1 | grep .proto$ ); do
      if [[ $protof =~ $protoReg ]]; then
        #echo "service $other from $service proto: ${BASH_REMATCH[1]}"
        #if [[ $other = $service ]]; then
        #  out=`dirname $protof`
        #else
        out="$ROOT"/services/$service/grpc/$other/${BASH_REMATCH[1]}.pb #`dirname $protof`
        #fi
        mkdir -p `dirname ./$out`
        pbjs $protof -t commonjs > $out.js
        #cp $protof $out
        #protoc --go_out=plugins=grpc:$out $protof
      fi
    done
  done
done

cd "$ROOT"


for service in $services; do
  cd "$ROOT"/services/$service
  docker build --no-cache -t $imagename .
done

