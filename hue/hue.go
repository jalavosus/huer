package hue

import (
	"context"
	"fmt"
	"time"

	"github.com/amimof/huego"

	"github.com/jalavosus/huer/entities"
	"github.com/jalavosus/huer/internal/config"
)

const tokenFilename string = "huetoken"

type Huer struct {
	uri      string
	username string
	bridge   *huego.Bridge
	Rooms    []*entities.Room
}

func NewHuer(uri, username string) (*Huer, error) {
	h := new(Huer)
	h.uri = uri

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	bridge, err := huego.DiscoverContext(ctx)
	if err != nil {
		return nil, err
	}
	h.bridge = bridge

	if username != "" {
		h.bridge = h.bridge.Login(username)
		h.username = username

		return h, nil
	}

	userToken, err := loadUserToken()
	if err != nil {
		return nil, err
	}

	if userToken == "" {
		return nil, fmt.Errorf("no token")
	}

	h.username = userToken
	h.bridge = h.bridge.Login(h.username)

	return h, nil
}

func (h *Huer) Bridge() *huego.Bridge {
	return h.bridge
}

func (h *Huer) AddRoom(room *entities.Room) {
	h.Rooms = append(h.Rooms, room)
}

func (h *Huer) AllLights() error {
	lights, err := h.bridge.GetLights()
	if err != nil {
		return err
	}

	for _, l := range lights {
		fmt.Println(l.Name, l.ID, l.State.On)
	}

	return nil
}

func (h *Huer) ToggleLight(lightId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultContextTimeout)
	defer cancel()

	light, err := h.bridge.GetLightContext(ctx, lightId)
	if err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(context.Background(), config.DefaultContextTimeout)
	defer cancel()

	_, err = h.bridge.SetLightStateContext(ctx, lightId, huego.State{On: !light.State.On})
	return err
}