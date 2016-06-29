#!/bin/bash
set -e
set -v
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..
cd "$ROOT"

./deploy/go/build.sh saver avezila/saver
./deploy/go/build.sh proxy avezila/saver-proxy
./deploy/mongo/build.sh saver-mongo avezila/saver-mongo
./deploy/node/build.sh saver-client avezila/saver-client

rm -rf build
