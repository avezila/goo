CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o ./app .

docker build -t $1 .

rm ./app
