package entities

import (
	"context"
	"fmt"

	"github.com/amimof/huego"

	"github.com/jalavosus/huer/utils"
)

type Light struct {
	*BaseEntity
	MacAddress string
	state      EntityState
	light      *huego.Light
}

func NewLight(l *huego.Light) *Light {
	newLight := &Light{
		BaseEntity: &BaseEntity{},
		MacAddress: l.UniqueID,
		light:      l,
	}

	newLight.SetName(l.Name)
	newLight.SetId(l.ID)
	newLight.SetUid(utils.MakeEntityUID(l.Name, l.ID, l.UniqueID))

	return newLight
}

func NewLightFromOpts(opts ...BaseEntityOpt) *Light {
	return &Light{
		BaseEntity: NewBaseEntityFromOpts(opts...),
	}
}

func (l Light) State() EntityState {
	return l.state
}

func (l Light) IsOn() bool {
	return l.state == StateOn
}

func (l Light) IsOff() bool {
	return l.state == StateOff
}

func (l Light) checkLightState() EntityState {
	if l.light.IsOn() {
		return StateOn
	}

	return StateOff
}

func (l *Light) Toggle() error {
	return utils.WithTimeoutCtx(func(ctx context.Context) error {
		switch l.IsOn() {
		case true:
			if err := l.light.OffContext(ctx); err != nil {
				return err
			}
			l.state = StateOff
		case false:
			if err := l.light.OnContext(ctx); err != nil {
				return err
			}
			l.state = StateOn
		}

		return fmt.Errorf("wat")
	})
}

func (l *Light) ToggleOn() error {
	if l.checkLightState() == StateOn {
		return nil
	}

	return l.Toggle()
}

func (l *Light) ToggleOff() error {
	if l.checkLightState() == StateOff {
		return nil
	}

	return l.Toggle()
}

func (l *Light) Rename(newName string) error {
	return utils.WithTimeoutCtx(func(ctx context.Context) error {
		return l.light.RenameContext(ctx, newName)
	})
}