package hue

import (
	"context"
	"fmt"

	"github.com/amimof/huego"

	"github.com/jalavosus/huer/entities"
	"github.com/jalavosus/huer/utils"

	"github.com/pkg/errors"
)

type Huer struct {
	uri       string
	userToken string
	bridge    *huego.Bridge
	Rooms     []*entities.Room
}

func NewHuer(uri, newUsername string) (*Huer, error) {
	h, err := newHuer(uri)
	if err != nil {
		return nil, err
	}

	h.userToken, err = h.createUser(newUsername)
	if err != nil {
		return nil, err
	}

	_, h.bridge = h.tryLogin(h.userToken)

	return h, nil
}

func NewHuerWithToken(uri, userToken string) (*Huer, error) {
	if userToken == "" {
		return nil, errors.New("param userToken can't be empty")
	}

	h, err := newHuer(uri)
	if err != nil {
		return nil, err
	}

	var authorized bool
	authorized, h.bridge = h.tryLogin(userToken)
	if !authorized {
		return nil, errors.Errorf("token %[1]s unauthorized", userToken)
	}

	h.userToken = userToken

	return h, nil
}

func newHuer(uri string) (*Huer, error) {
	if uri == "" {
		return nil, fmt.Errorf("no uri provided")
	}

	h := new(Huer)
	h.uri = uri

	return h, utils.WithTimeoutCtx(func(ctx context.Context) error {
		bridge, err := huego.DiscoverContext(ctx)
		if err != nil {
			return err
		}

		h.bridge = bridge

		return nil
	})
}

func (h *Huer) Bridge() *huego.Bridge {
	return h.bridge
}

func (h Huer) Token() string {
	return h.userToken
}