#!/bin/bash
set -e
set -v
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/../..
cd "$ROOT"

if [[ $# < 2 ]]; then
	echo "usage:  docker.sh service-name image-name"
	exit 1
fi


imagename=$2
servicename=$1
appname=`echo $imagename | tr "/" __` 
appname=`echo $appname | tr ":" _`

rm -rf build/app
mkdir -p build/app
cp -r services/$servicename build/app/mongo
cp deploy/mongo/Dockerfile build/app
cd build/app
docker build --no-cache -t $imagename .

cd "$ROOT"
rm -rf build/app