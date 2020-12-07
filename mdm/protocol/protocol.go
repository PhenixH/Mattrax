package protocol

import (
	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
)

type Protocol interface {
	ID() string
	Init(srv *mattrax.Server) error
	Mount() error
	MountAPI(r *mux.Router, rUnauthenticated *mux.Router) error
	Events() EventHandlers
	Status() (interface{}, error)
}
