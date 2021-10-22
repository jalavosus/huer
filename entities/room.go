package entities

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jalavosus/huer/internal/config"
	"github.com/jalavosus/huer/internal/params"
)

type Room struct {
	*Entity `yaml:",inline"`
	Lights  []*Light `json:"lights" yaml:"lights"`
}

func (r *Room) ID(h Huer) int {
	if r.Entity.ID != 0 {
		return r.Entity.ID
	}

	if h != nil {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), config.DefaultContextTimeout)
			defer cancel()

			grps, err := h.Bridge().GetGroupsContext(ctx)
			if err != nil {
				log.Fatalln(err)
			}

			for _, grp := range grps {
				if strings.ToLower(grp.Name) == strings.ToLower(r.Name) {
					r.Entity.ID = grp.ID
					return
				}
			}
		}()
	}

	return r.Entity.ID
}

func (r *Room) LightsInfo(h Huer) ([]*Light, error) {
	var lights []*Light

	if len(r.Lights) == 0 {
		grp, err := h.Bridge().GetGroup(r.ID(h))
		if err != nil {
			return nil, err
		}

		for _, l := range grp.Lights {
			id, _ := strconv.ParseInt(l, 10, 32)
			r.Lights = append(r.Lights, &Light{
				ID: int(id),
			})
		}
	}

	for _, l := range r.Lights {
		args := params.NewRoomArgs(params.NameParam(l.Name), params.IDParam(l.ID))

		light, err := r.Light(h, args)
		if err != nil {
			return nil, err
		}

		lights = append(lights, light)
	}

	return lights, nil
}

func (r *Room) Light(h Huer, args *params.RoomArgs) (*Light, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultContextTimeout)
	defer cancel()

	switch {
	case args.HasID():
		light, err := h.Bridge().GetLightContext(ctx, args.ID())
		if err != nil {
			return nil, err
		}

		return NewLight(light), nil
	case args.HasName():
		lights, err := h.Bridge().GetLightsContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, l := range lights {
			if strings.ToLower(l.Name) == strings.ToLower(args.Name()) {
				return NewLight(&l), nil
			}
		}
	}

	return nil, fmt.Errorf("can't get light info")
}