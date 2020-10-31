package server

import (
	"log"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/asukhodko/go-grps-cache-and-consumer/pkg/proto/randomdatastream"
)

type Server interface {
	Serve() error
}

func NewServer(port string) Server {
	s := &server{
		port:       port,
		grpcServer: grpc.NewServer(),
	}
	reflection.Register(s.grpcServer)
	pb.RegisterRandomDataStreamerServer(s.grpcServer, s)
	return s
}

type server struct {
	pb.UnimplementedRandomDataStreamerServer
	port       string
	grpcServer *grpc.Server
}

func (s *server) Serve() error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		return errors.Wrap(err, "s.grpcServer.Serve")
	}

	return nil
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
