package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/jalavosus/huer/internal/config"
	"github.com/jalavosus/huer/server/magic"
	"github.com/jalavosus/huer/server/serverutils"
)

const (
	DefaultMagicHeader string = "X-Huer-Magic-Fuckery"
	DefaultPort        string = "44562"
)

const readWriteTimeout = 15 * time.Second

type Server struct {
	conf *config.Config
}

func NewServer(conf *config.Config) *Server {
	s := &Server{conf}

	if s.conf.MagicHeader.Header == "" {
		s.conf.MagicHeader.Header = DefaultMagicHeader
	}

	return s
}

func (s *Server) Start() error {
	r := mux.NewRouter()
	r.Use(magic.Middleware(s.conf.MagicHeader))

	r.HandleFunc("/",
		serverutils.SimpleMessageHandler(
			"You made it!",
			http.StatusOK),
	).
		Methods("GET")

	sr := r.PathPrefix("/api").
		Subrouter().
		StrictSlash(true)

	sr.HandleFunc("/",
		serverutils.SimpleMessageHandler(
			"It's the API's root path, the hell do you want?",
			http.StatusOK,
		),
	).
		Methods("GET")

	makeBasicToggleable(sr, "/rooms")

	makeBasicToggleable(sr, "/lights")

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:" + DefaultPort,
		WriteTimeout: readWriteTimeout,
		ReadTimeout:  readWriteTimeout,
	}

	log.Println("server listening on " + srv.Addr)

	return srv.ListenAndServe()
}