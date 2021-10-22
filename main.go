package main

import (
	"fmt"
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

	fmt.Println(conf)

	s := server.NewServer(conf)

	log.Fatalln(s.Start())
}