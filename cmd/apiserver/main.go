package main

import (
	"flag"
	"fmt"
	"github.com/VSKrivoshein/short-link/internal/app/apiserver"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.yaml", "path to config file")
}

func main() {
	flag.Parse()
	fmt.Println("configPath", configPath)

	if err := apiserver.InitConfig(configPath); err != nil {
		log.Fatalf("config is not initialized: %s", err.Error())
	}

	config := apiserver.NewConfig()

	s := apiserver.NewAPIServer(config)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
