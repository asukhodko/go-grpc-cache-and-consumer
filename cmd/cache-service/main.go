package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/asukhodko/go-grps-cache-and-consumer/pkg/cachingfetcher"
	"github.com/asukhodko/go-grps-cache-and-consumer/pkg/server"
	"github.com/asukhodko/go-grps-cache-and-consumer/pkg/service"
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

	f := cachingfetcher.NewFetcher()
	svc := service.NewService(f, config.URLs, config.MinTimeout, config.MaxTimeout, config.NumberOfRequests)
	srv := server.NewServer(port, svc)

	err = srv.Serve()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
