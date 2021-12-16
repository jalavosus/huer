package entities

import (
	"context"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/amimof/huego"
	"github.com/pkg/errors"

	"github.com/jalavosus/huer/utils"
)

type roomGroup struct {
	h    Huer
	id   int
	once sync.Once
	grp  *huego.Group
}

func newRoomGroup(h Huer, id int) *roomGroup {
	g := new(roomGroup)
	g.h = h
	g.id = id

	return g
}

func (g *roomGroup) init() {
	g.once.Do(func() {
		grp, err := g.h.Bridge().GetGroup(g.id)
		if err != nil {
			log.Fatalln(err)
		}

		g.grp = grp
	})
}

func (g *roomGroup) Group() *huego.Group {
	g.init()
	return g.grp
}

type Room struct {
	*BaseEntity `yaml:",inline"`
	Lights      []*Light `json:"lights" yaml:"lights"`
	idOnce      sync.Once
	grpOnce     sync.Once
	group       *roomGroup
}

func NewRoomFromOpts(opts ...BaseEntityOpt) *Room {
	r := &Room{
		BaseEntity: NewBaseEntityFromOpts(opts...),
	}

	if r.Hue() != nil {
		r.grpOnce.Do(func() {
			r.group = newRoomGroup(r.Huer(), r.Id())
			r.group.init()
		})
	}

	return r
}

func (r *Room) grp() *huego.Group {
	r.grpOnce.Do(func() {
		g := newRoomGroup(r.Huer(), r.Id())
		g.init()

		r.group = g
	})

	return r.group.Group()
}

func (r *Room) state() *huego.State {
	return r.grp().State
}

func (r *Room) RawColors() ([]float32, uint8) {
	s := r.state()

	return s.Xy, s.Bri
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

	r.idOnce.Do(func() {
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
		grp := r.grp()

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