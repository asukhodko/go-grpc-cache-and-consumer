package server

import (
	"log"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/asukhodko/go-grps-cache-and-consumer/pkg/proto/randomdatastream"
	"github.com/asukhodko/go-grps-cache-and-consumer/pkg/service"
)

// Server - интерфейс grpc-сервера
type Server interface {
	Serve() error
}

// NewServer конструирует сервер
func NewServer(port string, service service.Service) Server {
	s := &server{
		port:       port,
		grpcServer: grpc.NewServer(),
		service:    service,
	}
	reflection.Register(s.grpcServer)
	pb.RegisterRandomDataStreamerServer(s.grpcServer, s)
	return s
}

type server struct {
	pb.UnimplementedRandomDataStreamerServer
	port       string
	grpcServer *grpc.Server
	service    service.Service
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
	ctx := stream.Context()
	chData, chErr := s.service.GetDataWithChannel(ctx)

	var data []byte
	ok := true
	for {
		select {
		case data, ok = <-chData:
			if !ok {
				break
			}
			err := stream.Send(&pb.GetRandomDataStreamResponse{Data: data})
			if err != nil {
				return errors.Wrap(err, "stream.Send")
			}
		case err := <-chErr:
			err = errors.Wrap(err, "s.service.GetDataWithChannel")
			log.Printf("[DEBUG] error: %v\n", err)
			return err
		}
		if !ok {
			break
		}
	}

	return nil
}
