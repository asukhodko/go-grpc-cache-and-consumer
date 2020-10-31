package main

import (
	"log"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/asukhodko/go-grps-cache-and-consumer/pkg/proto/randomdatastream"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedRandomDataStreamerServer
}

func (s *server) GetRandomDataStream(_ *pb.GetRandomDataStreamRequest, stream pb.RandomDataStreamer_GetRandomDataStreamServer) error {
	log.Println("Invoked.")
	err := stream.Send(&pb.GetRandomDataStreamResponse{Data: "message 1"})
	if err != nil {
		return errors.Wrap(err, "stream.Send")
	}
	err = stream.Send(&pb.GetRandomDataStreamResponse{Data: "message 2"})
	if err != nil {
		return errors.Wrap(err, "stream.Send")
	}
	err = stream.Send(&pb.GetRandomDataStreamResponse{Data: "message 3"})
	if err != nil {
		return errors.Wrap(err, "stream.Send")
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterRandomDataStreamerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
