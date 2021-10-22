package entities

import (
	"context"
	"fmt"

	"github.com/amimof/huego"

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
	return !l.IsOn()
}

func (l *Light) Toggle() error {
	return utils.WithTimeoutCtx(func(ctx context.Context) error {
		switch l.IsOn() {
		case true:
			return l.light.OffContext(ctx)
		case false:
			return l.light.OnContext(ctx)
		}

		return fmt.Errorf("wat")
	})
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
	return utils.WithTimeoutCtx(func(ctx context.Context) error {
		return l.light.RenameContext(ctx, newName)
	})
}