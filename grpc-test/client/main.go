package main

import (
	"context"
	"exercise/grpc-test/client/proto/hello"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial(":33333", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := hello.NewHelloClient(conn)

	res, err := c.SayHello(context.Background(), &hello.HelloRequest{Name: "hello world"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.Message)

	stream, err := c.LotsOfReplies(context.Background(), &hello.HelloRequest{Name: "test"})
	if err != nil {
		log.Fatal(err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("stream recv:", err)
		}
		log.Println(res.Message)
	}
}
