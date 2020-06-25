package main

import (
	"flag"
	"github.com/Alex27Khalupka/Go-course-task/pkg/apiserver"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()

	// getting config
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	// Server configuration
	s := apiserver.New(config)

	// starting server
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
