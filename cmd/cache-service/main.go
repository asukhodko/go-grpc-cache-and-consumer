package main

import (
	"log"

	"github.com/asukhodko/go-grps-cache-and-consumer/pkg/server"
)

const (
	port = ":50051"
)

func main() {
	s := server.NewServer(port)

	err := s.Serve()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
