package main

import (
	"flag"
	"github.com/dantedenis/reast-api-golang/internal/app/apiserver"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.yml", "path to config api server file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(file)

	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Fatalln(err)
	}

	s := apiserver.NewAPIServer(config)
	if err := s.Start(); err != nil {
		log.Fatalln(err)
	}

}
