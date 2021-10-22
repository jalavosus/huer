package magic

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/jalavosus/huer/internal/config"
	"github.com/jalavosus/huer/server/serverutils"
)

func Middleware(conf *config.MagicHeaderConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
				return
			}

			magicHeaderVal := r.Header.Get(conf.Header)

			magicValue, err := Decrypt(conf.Key, magicHeaderVal)
			if err != nil {
				log.Println(err)
				invalidMagicResponse(true, w)
				return
			} else if magicValue != conf.Value {
				invalidMagicResponse(false, w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func invalidMagicResponse(isErr bool, w http.ResponseWriter) {
	var (
		statusCode = http.StatusForbidden
		msg        = "You used the wrong kind of magic, lad"
	)
	if isErr {
		statusCode = http.StatusInternalServerError
		msg = "You used some real trippy magic, bruh"
	}

	serverutils.SendSimpleResponseWithStatusCode(msg, statusCode, w)
	return
}