package entities

import (
	"github.com/amimof/huego"
)

type Entity struct {
	Name string `json:"name" yaml:"name"`
	ID   int    `json:"id" yaml:"id"`
	UID  string `json:"-" yaml:"-"`
}

type Huer interface {
	Bridge() *huego.Bridge
}