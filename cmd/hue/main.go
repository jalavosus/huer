package main

import (
	"fmt"
	"log"

	"github.com/jalavosus/huer/entities"
	"github.com/jalavosus/huer/hue"
	"github.com/jalavosus/huer/internal/config"
)

const (
	configFile string = "./huer.yaml"
)

func main() {
	conf, err := config.LoadConfig(configFile)
	if err != nil {
		log.Panic(err)
	}

	h, err := hue.NewHuerWithToken(conf.URI, conf.Token)
	if err != nil {
		log.Println(err)
		return
	}

	h.AddRoom(&entities.Room{
		BaseEntity: entities.NewBaseEntityFromOpts(entities.EntityName("Bedroom")),
	})

	// rooms, err := h.LoadRooms()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for _, room := range h.Rooms {
		fmt.Println("==== " + room.Name() + " " + room.Uid() + " ====")

		lights, err := room.LightsInfo()
		if err != nil {
			log.Fatal(err)
		}

		for _, light := range lights {
			fmt.Println(light.Name() + "\t" + light.Uid())
		}

		fmt.Println()
	}

}