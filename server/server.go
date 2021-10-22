package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/jalavosus/huer/internal/config"
	"github.com/jalavosus/huer/server/magic"
	"github.com/jalavosus/huer/server/serverutils"
)

const DefaultMagicHeader string = "HUER_MAGIC_FUCKERY"

type Server struct {
	mux  *mux.Router
	conf *config.Config
}

func NewServer(conf *config.Config) *Server {
	s := &Server{
		conf: conf,
	}

	if s.conf.MagicHeader.Header == "" {
		s.conf.MagicHeader.Header = DefaultMagicHeader
	}

	return s
}

func (s *Server) Start() error {
	r := mux.NewRouter()
	r.Use(magic.Middleware(s.conf.MagicHeader))

	r.HandleFunc("/", serverutils.SimpleMessageHandler("You made it!", http.StatusOK))

	sr := r.PathPrefix("/api").
		Subrouter().
		StrictSlash(true)

	sr.HandleFunc("/", serverutils.SimpleMessageHandler("It's the API's root path, the hell do you want?", http.StatusOK)).
		Methods("GET")

	makeBasicToggleable(sr, "/rooms")

	makeBasicToggleable(sr, "/lights")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:44562",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return srv.ListenAndServe()
}

func makeBasicToggleable(r *mux.Router, pathPrefix string) {
	if !strings.HasPrefix(pathPrefix, "/") {
		pathPrefix = "/" + pathPrefix
	}

	if strings.HasSuffix(pathPrefix, "/") {
		pathPrefix = strings.TrimSuffix(pathPrefix, "/")
	}

	sr := r.PathPrefix(pathPrefix).
		Subrouter().
		StrictSlash(true)

	sr.HandleFunc("/", nil).
		Methods("GET")

	sr.HandleFunc("/{name}", nil).
		Methods("GET")

	sr.HandleFunc("/{name}/toggle", nil).
		Methods("POST")

	sr.HandleFunc("/{name}/toggle/{state}", nil).
		Methods("POST")
}