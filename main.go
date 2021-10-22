package main

import (
	"log"

	"github.com/jalavosus/huer/internal/config"
	"github.com/jalavosus/huer/server"
)

const (
	configFile string = "./huer.yaml"
)

func main() {
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Panic(err)
	}

	s := server.NewServer(conf)

	log.Fatalln(s.Start())
}