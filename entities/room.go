package entities

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/jalavosus/huer/utils"
)

type Room struct {
	*BaseEntity `yaml:",inline"`
	Lights      []*Light `json:"lights" yaml:"lights"`
	once        sync.Once
}

func (r *Room) Id() int {
	id, err := r.getId(r.BaseEntity.Id())
	if err != nil {
		log.Fatalln(err)
	}

	return id
}

func (r *Room) getId(roomId int) (int, error) {
	if roomId != -1 {
		return roomId, nil
	}

	var (
		foundId = -1
		idErr   error
	)

	r.once.Do(func() {
		if r.Hue() == nil {
			idErr = errors.Errorf("no huego.Bridge set for object")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		grps, err := r.Hue().GetGroupsContext(ctx)
		if err != nil {
			idErr = err
			return
		}

		for _, grp := range grps {
			if strings.ToLower(grp.Name) == strings.ToLower(r.Name()) {
				foundId = grp.ID
				break
			}
		}
	})

	if foundId == -1 {
		idErr = errors.Errorf("couldn't get the room ID for room %[1]s", r.Name())
	} else if foundId != -1 && idErr == nil {
		r.SetId(foundId)
	}

	return foundId, idErr
}

func (r *Room) LightsInfo() ([]*Light, error) {
	if r.Hue() == nil {
		return nil, errors.Errorf("no huego.Bridge set for object")
	}

	var lights []*Light

	if len(r.Lights) == 0 {
		grp, err := r.Hue().GetGroup(r.Id())
		if err != nil {
			return nil, err
		}

		for _, l := range grp.Lights {
			id, _ := strconv.ParseInt(l, 10, 32)
			r.Lights = append(r.Lights, NewLightFromOpts(EntityId(int(id))))
		}
	}

	for _, l := range r.Lights {
		args := []BaseEntityOpt{EntityId(l.Id())}
		light, err := r.Light(args...)
		if err != nil {
			return nil, err
		}

		lights = append(lights, light)
	}

	return lights, nil
}

func (r *Room) Light(opts ...BaseEntityOpt) (*Light, error) {
	if r.Hue() == nil {
		return nil, errors.Errorf("no huego.Bridge set for object")
	}

	var res *Light

	base := NewBaseEntityFromOpts(opts...)

	err := utils.WithTimeoutCtx(func(ctx context.Context) error {
		switch {
		case base.Id() != -1:
			l, err := r.Hue().GetLightContext(ctx, base.Id())
			if err != nil {
				return err
			}

			res = NewLight(l)
		case base.Name() != "":
			lights, err := r.Hue().GetLightsContext(ctx)
			if err != nil {
				return err
			}

			for _, l := range lights {
				if strings.ToLower(l.Name) == strings.ToLower(base.Name()) {
					res = NewLight(&l)
					break
				}
			}
		default:
			return nil
		}

		return nil
	})

	return res, err
}