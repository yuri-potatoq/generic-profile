package server

import (
	"github.com/gorilla/mux"
	v1 "github.com/yuri-potatoq/generic-profile/api/v1"
	"net/http"
)

type ServerOpts struct {
	port int

	GetEnrollmentHandler   *v1.GetEnrollmentHandler
	PatchEnrollmentHandler *v1.PatchEnrollmentHandler
}

type Server struct {
	opts *ServerOpts
}

func NewServer(opts *ServerOpts) *Server {
	return &Server{opts}
}

func (s *Server) Start() error {

	r := mux.NewRouter()
	r.Handle("/profile/{id:[0-9]+}", s.opts.GetEnrollmentHandler).Methods(http.MethodGet)
	r.Handle("/profile", s.opts.PatchEnrollmentHandler).Methods(http.MethodPatch)

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	return srv.ListenAndServe()
}
