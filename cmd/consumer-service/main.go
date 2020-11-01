package main

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"

	pb "github.com/asukhodko/go-grps-cache-and-consumer/pkg/proto/randomdatastream"
)

const (
	address        = "localhost:50051"
	timeoutSeconds = 2
	routines       = 10
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRandomDataStreamerClient(conn)

	wg := sync.WaitGroup{}
	wg.Add(routines)
	for id := 1; id <= routines; id++ {
		go func(id int) {
			defer wg.Done()
			receiveData(c, id)
		}(id)
	}
	wg.Wait()
}

func receiveData(c pb.RandomDataStreamerClient, id int) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds*time.Second)
	defer cancel()

	r, err := c.GetRandomDataStream(ctx, &pb.GetRandomDataStreamRequest{})
	if err != nil {
		log.Printf("[%d] could not call GetRandomDataStream: %v\n", id, err)
		return
	}
	for n := 1; ; n++ {
		msg, err := r.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("[%d] r.Recv: %v\n", id, err)
			return
		}
		if msg == nil {
			log.Printf("[%d] response message is nil\n", id)
			return
		}
		log.Printf("[%d] Received stream %d (%d bytes).\n", id, n, len(msg.Data))
	}
}
