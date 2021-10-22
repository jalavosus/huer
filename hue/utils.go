package hue

import (
	"context"
	"io/ioutil"
	"time"
)

func (h *Huer) createUser(username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	time.Sleep(30 * time.Second)
	return h.bridge.CreateUserContext(ctx, username)
}

func loadUserToken() (string, error) {
	data, err := ioutil.ReadFile(tokenFilename)
	if err != nil {
		return "", err
	}

	return string(data), nil
}