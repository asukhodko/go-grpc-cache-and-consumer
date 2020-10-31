package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb "github.com/asukhodko/go-grps-cache-and-consumer/pkg/proto/randomdatastream"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRandomDataStreamerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetRandomDataStream(ctx, &pb.GetRandomDataStreamRequest{})
	if err != nil {
		log.Fatalf("could not call GetRandomDataStream: %v", err)
	}
	for {
		msg, err := r.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			err = errors.Wrap(err, "r.Recv")
			fmt.Println(err)
			return
		}
		if msg == nil {
			err = errors.New("response message is nil")
			fmt.Println(err)
			return
		}
		fmt.Printf("Received: %s\n", msg.Data)
	}
	log.Printf("Response: %v", r)
}
