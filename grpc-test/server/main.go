package main

import (
	hello2 "exercise/grpc-test/server/controller/hello"
	"exercise/grpc-test/server/proto/hello"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":33333")
	if err != nil {
		log.Fatal("fail listen :", err)
	}

	s := grpc.NewServer()
	hello.RegisterHelloServer(s, &hello2.Controller{})
	log.Println("listen on port 33333")
	if err = s.Serve(listen); err != nil {
		log.Fatal("fail serve :", err)
	}
}
