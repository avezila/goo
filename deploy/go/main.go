package main

import (
	"net"
	"os"

	"google.golang.org/grpc"
)

type server struct {
	listener net.Listener
	grpc     *grpc.Server
}

func newServer(port string) (serv *server, err error) {
	serv = &server{}
	serv.listener, err = net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}
	serv.grpc = grpc.NewServer()
	return serv, nil
}

func (serv *server) serve() {
	serv.grpc.Serve(serv.listener)
}

func main() {

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "5353"
	}
	serv, err := newServer(":" + port)

	if err != nil {
		panic(err)
	}

	if err = register(serv); err != nil {
		panic(err)
	}

	serv.serve()
}
