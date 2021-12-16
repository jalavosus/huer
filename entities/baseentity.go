package entities

import (
	"encoding/json"

	"github.com/amimof/huego"

	"github.com/jalavosus/huer/utils"
)

type BaseEntity struct {
	name string
	id   int
	uid  string
	hue  Huer
}

func newDefaultBaseEntity() *BaseEntity {
	return &BaseEntity{
		id: -1,
	}
}

func NewBaseEntityFromOpts(opts ...BaseEntityOpt) *BaseEntity {
	e := newDefaultBaseEntity()
	ApplyOpts(e, opts...)

	return e
}

func (e BaseEntity) Name() string {
	return e.name
}

func (e *BaseEntity) SetName(name string) {
	e.name = name
}

func (e BaseEntity) Id() int {
	return e.id
}

func (e *BaseEntity) SetId(id int) {
	e.id = id
}

func (e BaseEntity) Uid() string {
	return e.uid
}

func (e *BaseEntity) SetUid(uid string) {
	e.uid = uid
}

func (e BaseEntity) Hue() *huego.Bridge {
	if e.hue != nil {
		return e.hue.Bridge()
	}

	return nil
}

func (e *BaseEntity) SetHue(huer Huer) {
	e.hue = huer
}

func (e BaseEntity) Huer() Huer {
	return e.hue
}

func (e *BaseEntity) SetHuer(huer Huer) {
	e.SetHue(huer)
}

func (e BaseEntity) marshal(marshal utils.MarshalFunc) ([]byte, error) {
	raw := make(map[string]interface{})

	if e.name != "" {
		raw["name"] = e.name
	}

	if e.id != -1 {
		raw["id"] = e.id
	}

	if e.uid != "" {
		raw["uid"] = e.uid
	}

	return marshal(raw)
}

func (e *BaseEntity) unmarshal(data []byte, unmarshal utils.UnmarshalFunc) error {
	var raw map[string]interface{}

	if err := unmarshal(data, &raw); err != nil {
		return err
	}

	if name := rawMapValue[string](raw, "name"); name != "" {
		e.name = name
	}

	if id := rawMapValue[int](raw, "id"); id != 0 {
		e.id = id
	}

	if uid := rawMapValue[string](raw, "uid"); uid != "" {
		e.uid = uid
	}

	return nil
}

func (e BaseEntity) MarshalJSON() ([]byte, error) {
	return e.marshal(json.Marshal)
}

func (e *BaseEntity) UnmarshalJSON(data []byte) error {
	return e.unmarshal(data, json.Unmarshal)
}

func (e BaseEntity) MarshalYAML() ([]byte, error) {
	return e.marshal(utils.MarshalYaml)
}

func (e *BaseEntity) UnmarshalYAML(data []byte) error {
	return e.unmarshal(data, utils.UnmarshalYAML)
}

func rawMapValue[T any](raw map[string]interface{}, key string) T {
	var defaultVal T

	if rawVal, ok := raw[key]; ok {
		if val, ok := rawVal.(T); ok {
			return val
		}
	}

	return defaultVal
}