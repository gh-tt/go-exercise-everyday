package hello

import (
	"context"
	"exercise/grpc-test/server/proto/hello"
	"fmt"
)

type Controller struct{}

func (h *Controller) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	return &hello.HelloResponse{Message: fmt.Sprintf("%s", in.Name)}, nil
}

func (h *Controller) LotsOfReplies(in *hello.HelloRequest, stream hello.Hello_LotsOfRepliesServer) error {
	for i := 0; i < 10; i++ {
		stream.Send(&hello.HelloResponse{Message: fmt.Sprintf("%s %s %d", in.Name, "Reply", i)})
	}
	return nil
}
