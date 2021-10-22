package hue

import (
	"context"
	"log"
	"time"

	"github.com/amimof/huego"

	"github.com/jalavosus/huer/utils"
)

func (h *Huer) tryLogin(userToken string) (authorized bool, authedBridge *huego.Bridge) {
	authedBridge = h.bridge.Login(userToken)

	authorized = utils.WithTimeoutCtx(func(ctx context.Context) error {
		_, err := authedBridge.GetConfigContext(ctx)
		if err != nil {
			return err
		}

		return nil
	}) == nil // gettin' real lazy with it

	return
}

func (h *Huer) createUser(username string) (userToken string, err error) {
	_ = utils.WithTimeoutCtx(func(ctx context.Context) error {
		log.Println("Just gonna wait ~30 seconds for you to hit the stupid button on your hue bridge")
		time.Sleep(29 * time.Second)

		userToken, err = h.bridge.CreateUserContext(ctx, username)
		return nil
	})

	return
}