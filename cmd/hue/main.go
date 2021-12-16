package main

import (
	"fmt"
	"log"

	"github.com/jalavosus/huer/entities"
	"github.com/jalavosus/huer/hue"
	"github.com/jalavosus/huer/hue/hueconsts"
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

	log.Println(h.Token())

	h.AddRoom(entities.NewRoomFromOpts(
		entities.EntityName("Bedroom"),
		entities.EntityHuer(h),
	))

	for _, room := range h.Rooms {
		fmt.Println("==== " + room.Name() + " " + room.Uid() + " ====")

		lights, err := room.LightsInfo()
		if err != nil {
			log.Fatal(err)
		}

		newB := hueconsts.CalculateMaxBrightnessDeltaMul(7, 3)

		for _, light := range lights {
			fmt.Println(light.Name())

			fmt.Println(light.RawColor())

			if err := light.SetColorAndBrightness(hueconsts.Nightlight, newB); err != nil {
				log.Fatalln(err)
			}

			fmt.Println(light.RawColor())
		}

		fmt.Println()
	}

}