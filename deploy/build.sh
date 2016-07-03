#!/bin/bash
set -e
set -v
ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"/..
cd "$ROOT"

cd services

cd goo
./build.sh avezila/goo
cd ..

cd goo-mongo
docker build -t avezila/goo-mongo .
cd ..

cd goo-web
docker build -t avezila/goo-web .
cd "$ROOT"
