package hue

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amimof/huego"

	"github.com/jalavosus/huer/entities"
	"github.com/jalavosus/huer/internal/params"
	"github.com/jalavosus/huer/utils"
)

func (h *Huer) GetRoomsRaw() ([]huego.Group, error) {
	return h.bridge.GetGroups()
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

		var roomLights []*entities.Entity
		for _, l := range r.Lights {
			id, _ := strconv.ParseInt(l, 10, 32)
			roomLights = append(roomLights, &entities.Entity{ID: int(id)})
		}

		h.Rooms = append(h.Rooms, &entities.Room{
			Entity: &entities.Entity{
				Name: r.Name,
				ID:   r.ID,
				UID:  rmUid,
			},
			Lights: roomLights,
		})
	}

	return h.Rooms, nil
}

func (h *Huer) ToggleRoom(args *params.RoomArgs) error {
	var id = -1

	if !args.HasName() && !args.HasID() {
		return fmt.Errorf("no room name or room ID provided")
	}

	if !args.HasID() && args.HasName() {
		grps, err := h.GetRoomsRaw()
		if err != nil {
			return err
		}

		for _, grp := range grps {
			if strings.ToLower(grp.Name) == strings.ToLower(args.Name()) {
				id = grp.ID
				break
			}
		}

		if id == -1 {
			return fmt.Errorf("no room with name %[1]s found", args.Name())
		}
	} else if args.HasID() {
		id = args.ID()
	}

	g, err := h.bridge.GetGroup(id)
	if err != nil {
		return err
	}

	switch g.State.On {
	case true:
		return g.Off()
	default:
		return g.On()
	}
}

func (h *Huer) hasRoom(uid string) (exists bool) {
	for _, r := range h.Rooms {
		if r.UID == uid {
			exists = true
			break
		}
	}

	return
}