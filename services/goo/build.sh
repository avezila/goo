#CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o ./app .

go build -o ./app main.go

docker build --no-cache -t $1 .

rm ./app
