package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type ServerOpts struct {
	GetEnrollmentHandler http.HandlerFunc
}

type Server struct {
	opts *ServerOpts
}

func NewServer(opts *ServerOpts) *Server {
	return &Server{opts}
}

func (s *Server) Start() error {

	r := mux.NewRouter()
	r.HandleFunc("/profile/{id}", s.opts.GetEnrollmentHandler).Methods(http.MethodGet)

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	return srv.ListenAndServe()
}
