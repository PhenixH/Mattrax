package mdm

import (
	"github.com/mattrax/Mattrax/mdm/agent"
	"github.com/mattrax/Mattrax/mdm/windows"

	mattrax "github.com/mattrax/Mattrax/internal"
)

// Mount initialises each of the MDM protocols
func Mount(srv *mattrax.Server) {
	agent.Mount(srv)
	windows.Mount(srv)
}
