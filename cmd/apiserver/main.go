package main

import (
	"context"
	"flag"
	"github.com/dantedenis/reast-api-golang/internal/app/apiserver"
	conf "github.com/dantedenis/reast-api-golang/internal/app/config"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.json", "path to config api server file")
}

func main() {
	flag.Parse()
	config, err := conf.NewConfigBuilder().Parse(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	s := apiserver.NewAPIServer(config)
	if err := s.Start(context.Background()); err != nil {
		log.Fatalln(err)
	}
}
