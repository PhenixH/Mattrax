package protocol

import "github.com/mattrax/Mattrax/internal/db"

type EventHandlers struct {
	CreatePolicy func(p db.Policy) error
	UpdatePolicy func(p db.Policy) error
	DeletePolicy func(p db.Policy) error
}
