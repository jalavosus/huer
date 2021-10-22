package entities

import (
	"context"
	"fmt"

	"github.com/amimof/huego"

	"github.com/jalavosus/huer/internal/config"
	"github.com/jalavosus/huer/utils"
)

type Light struct {
	Name       string
	ID         int
	UID        string
	MacAddress string
	light      *huego.Light
}

func NewLight(l *huego.Light) *Light {
	return &Light{
		Name:       l.Name,
		ID:         l.ID,
		UID:        utils.MakeEntityUID(l.Name, l.ID, l.UniqueID),
		MacAddress: l.UniqueID,
		light:      l,
	}
}

func (l Light) IsOn() bool {
	return l.light.State.On
}

func (l Light) IsOff() bool {
	return !l.light.State.On
}

func (l *Light) Toggle() error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultContextTimeout)
	defer cancel()

	switch l.light.State.On {
	case true:
		return l.light.OffContext(ctx)
	case false:
		return l.light.OnContext(ctx)
	}

	return fmt.Errorf("wat")
}

func (l *Light) ToggleOn() error {
	if l.IsOn() {
		return nil
	}

	return l.Toggle()
}

func (l *Light) ToggleOff() error {
	if l.IsOff() {
		return nil
	}

	return l.Toggle()
}

func (l *Light) Rename(newName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultContextTimeout)
	defer cancel()

	return l.light.RenameContext(ctx, newName)
}