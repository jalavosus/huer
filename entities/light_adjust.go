package entities

import (
	"context"
	"time"

	"github.com/amimof/huego"

	"github.com/jalavosus/huer/hue/hueconsts"
)

const (
	maxBrightness uint8 = 254
)

// func (l Light) Color() (color.Color, error) {
//
// }

func (l Light) state() *huego.State {
	return l.light.State
}

func (l *Light) refresh() {
	_ = l.light.IsOn()
}

func (l *Light) RawColor() ([]float32, uint8) {
	return l.state().Xy, l.state().Bri
}

func (l *Light) SetColor(c hueconsts.Color) error {
	return l.setColorAndBrightness(c, l.CurrentBrightness())
}

func (l *Light) SetColorAndBrightness(c hueconsts.Color, b uint8) error {
	return l.setColorAndBrightness(c, b)
}

func (l *Light) setColorAndBrightness(c hueconsts.Color, b uint8) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	update := huego.State{
		On:  true,
		Xy:  c,
		Bri: b,
	}

	return l.light.SetStateContext(ctx, update)
}

func (l *Light) CurrentBrightness() uint8 {
	return l.state().Bri
}

func (l *Light) setBrightness(newBrightness uint8) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return l.light.BriContext(ctx, newBrightness)
}

func (l *Light) SetBrightness(newBrightness uint8) error {
	return l.setBrightness(newBrightness)
}

func (l *Light) MaxBrightness() error {
	return l.setBrightness(maxBrightness)
}

// func (l *Light) SetDaylight() error {
//
// }