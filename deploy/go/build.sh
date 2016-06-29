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

protoReg='(.*)\.proto$'



for service in $services; do
  rm -rf `find services/$service | grep "\.pb\.go$"`
  for other in $allservices; do
    cd "$ROOT"/services/$other
    for protof in $( ls -1 | grep .proto$ ); do
      if [[ $protof =~ $protoReg ]]; then
        #echo "service $other from $service proto: ${BASH_REMATCH[1]}"
        if [[ $other = $service ]]; then
          out=`dirname ./$protof`
        else
          out="$ROOT"/services/$service/grpc/$other/`dirname ./$protof`
        fi
        mkdir -p $out
        protoc --go_out=plugins=grpc:$out ./$protof
      fi
    done
  done
done

cd "$ROOT"

rm -rf build/$output.pkg
mkdir -p build/$output.pkg
cd build/$output.pkg



cat > register.go << EOF
package main

EOF


for service in $services; do
cat >> register.go << EOF
import "../../services/$service"
EOF
done

cat >> register.go << EOF

func register (serv * server) error {

EOF


for service in $services; do
cat >> register.go << EOF
	if _,err := $service.New(serv.grpc); err != nil {
		return err
	}
EOF
done

cat >> register.go << EOF

	return nil
}

EOF

cp "$ROOT"/deploy/go/main.go ./

CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o ../$output .

cd "$ROOT"
rm -rf build/$output.pkg

rm -rf build/app
mkdir -p build/app
cp build/$output build/app/app
cp deploy/go/Dockerfile build/app
cd build/app
docker build --no-cache -t $imagename .

cd "$ROOT"
rm -rf build/app build/$output
