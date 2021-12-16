package hue

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/amimof/huego"

	"github.com/jalavosus/huer/entities"
	"github.com/jalavosus/huer/utils"
)

func (h *Huer) addToEntities(hueEntities ...entities.HueEntity) {
	for i := range hueEntities {
		if hueEntities[i].Hue() == nil {
			hueEntities[i].SetHue(h)
		}
	}
}

func (h *Huer) AddRoom(room *entities.Room) {
	h.addToEntities(room)

	if len(room.Lights) == 0 {
		func() {
			rmLights, err := room.LightsInfo()
			if err != nil {
				log.Fatalln(err)
			}

			for i := range rmLights {
				rmLights[i].SetHue(h)
			}

			room.Lights = rmLights
		}()
	}

	h.Rooms = append(h.Rooms, room)
}

func (h *Huer) GetRoomsRaw() (grp []huego.Group, err error) {
	_ = utils.WithTimeoutCtx(func(ctx context.Context) error {
		grp, err = h.bridge.GetGroupsContext(ctx)
		return nil
	})

	return
}

func (h *Huer) LoadRooms() ([]*entities.Room, error) {
	rawRooms, err := h.GetRoomsRaw()
	if err != nil {
		return nil, err
	}

	for _, r := range rawRooms {
		rmUid := utils.MakeEntityUID(r.Name, r.ID, r.Type, r.Class)
		if h.hasRoom(rmUid) {
			continue
		}

		var roomLights []*entities.Light
		for _, l := range r.Lights {
			id, _ := strconv.ParseInt(l, 10, 32)
			roomLights = append(roomLights, entities.NewLightFromOpts(entities.EntityId(int(id))))
		}

		newRoom := &entities.Room{
			BaseEntity: entities.NewBaseEntityFromOpts(
				entities.EntityName(r.Name),
				entities.EntityId(r.ID),
				entities.EntityUid(rmUid),
			),
			Lights: roomLights,
		}

		h.addToEntities(newRoom)

		h.Rooms = append(h.Rooms, newRoom)
	}

	return h.Rooms, nil
}

func (h *Huer) ToggleRoom(args ...entities.BaseEntityOpt) error {
	var id = -1

	base := entities.NewBaseEntityFromOpts(args...)

	if base.Id() == -1 && base.Name() == "" {
		return fmt.Errorf("no room name or room ID provided")
	}

	if base.Id() == -1 && base.Name() != "" {
		grps, err := h.GetRoomsRaw()
		if err != nil {
			return err
		}

		for _, grp := range grps {
			if strings.ToLower(grp.Name) == strings.ToLower(base.Name()) {
				id = grp.ID
				break
			}
		}

		if id == -1 {
			return fmt.Errorf("no room with name %[1]s found", base.Name())
		}
	} else if base.Id() != -1 {
		id = base.Id()
	}

	return utils.WithTimeoutCtx(func(ctx context.Context) error {
		g, err := h.bridge.GetGroupContext(ctx, id)
		if err != nil {
			return err
		}

		switch g.State.On {
		case true:
			return g.Off()
		default:
			return g.On()
		}
	})
}

func (h *Huer) hasRoom(uid string) (exists bool) {
	for _, r := range h.Rooms {
		if r.Uid() == uid {
			exists = true
			break
		}
	}

	return
}