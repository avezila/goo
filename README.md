### Just another url shortener

##### Look at https://goo.avezila.com

##### For start
* Install docker-engine:
  https://docs.docker.com/engine/installation/linux/fedora/
* Install docker-compose: 
  https://docs.docker.com/compose/install/
* Then run ```cd deploy/localhost && docker-compose up```
* Now you can open http://localhost:3000/ and check

##### Build:
* install go compiler
* then run 
  ```./deploy/build.sh```  
  this will create local docker images avezila/goo avezila/goo-web avezila/goo-mongo
