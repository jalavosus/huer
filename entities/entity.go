package entities

import (
	"github.com/amimof/huego"
)

type HueEntity interface {
	Name() string
	Id() int
	Uid() string
	Hue() *huego.Bridge

	SetName(string)
	SetId(int)
	SetUid(string)
	SetHue(Huer)
}

type Huer interface {
	Bridge() *huego.Bridge
}