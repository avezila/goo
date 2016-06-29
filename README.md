### Test Saver service - for save bin data by id

#### For start
* Install docker-engine:
  https://docs.docker.com/engine/installation/linux/fedora/
* Install docker-compose: 
  https://docs.docker.com/compose/install/
* Then run ```./deploy/compose.sh up```
* Now you can open http://127.0.0.1:3000/ and check

#### Build:
* install node, go compiler, protobuf, grpc, go plugins for grpc, protobufjs from npm for pbjs
* then run ```./deploy/build-saver.sh```
* or make your compose config and others.
