package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type routePath struct {
	route, method string
	handler       http.HandlerFunc
}

var basicToggleableRoutePaths = []routePath{
	{"/", "GET", nil},
	{"/{name}", "GET", nil},
	{"/{name}/state", "GET", nil},
	{"/{name}/toggle", "POST", nil},
	{"/{name}/toggle/{state}", "POST", nil},
}

func makeBasicToggleable(r *mux.Router, pathPrefix string) {
	pathPrefix = strings.TrimSuffix(pathPrefix, "/")

	if !strings.HasPrefix(pathPrefix, "/") {
		pathPrefix = "/" + pathPrefix
	}

	sr := r.PathPrefix(pathPrefix).
		Subrouter().
		StrictSlash(true)

	for _, rp := range basicToggleableRoutePaths {
		sr.HandleFunc(rp.route, rp.handler).
			Methods(rp.method)
	}
}