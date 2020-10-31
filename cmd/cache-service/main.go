package main

import (
	"github.com/asukhodko/go-grps-cache-and-consumer/pkg/service"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/asukhodko/go-grps-cache-and-consumer/pkg/server"
)

const (
	port           = ":50051"
	configFilename = "config.yml"
)

func main() {

	configBody, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	config := &struct {
		URLs             []string `yaml:"URLs"`
		MinTimeout       int      `yaml:"MinTimeout"`
		MaxTimeout       int      `yaml:"MaxTimeout"`
		NumberOfRequests int      `yaml:"NumberOfRequests"`
	}{}
	err = yaml.Unmarshal(configBody, config)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	svc := service.NewService(config.URLs, config.MinTimeout, config.MaxTimeout, config.NumberOfRequests)
	srv := server.NewServer(port, svc)

	err = srv.Serve()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
