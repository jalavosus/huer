package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func MakeEntityUID(name string, id int, extra ...string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(strings.ToLower(name)))
	_, _ = h.Write([]byte(fmt.Sprintf("%d", id)))

	for _, e := range extra {
		_, _ = h.Write([]byte(e))
	}

	return hex.EncodeToString(h.Sum(nil))
}