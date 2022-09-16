package mdm

import (
	"github.com/mattrax/Mattrax/mdm/android"
	"github.com/mattrax/Mattrax/mdm/protocol"
	"github.com/mattrax/Mattrax/mdm/windows"
	"github.com/rs/zerolog/log"

	mattrax "github.com/mattrax/Mattrax/internal"
)

var Protocols []protocol.Protocol

// Mount initialises each of the MDM protocols
func Mount(srv *mattrax.Server) {
	if srv.Args.GoogleServiceAccountPath == "" {
		log.Warn().Msg("Google Android protocol disabled.")
	} else {
		Protocols = append(Protocols, &android.Protocol{})
	}
	Protocols = append(Protocols, &windows.Protocol{})

	for _, p := range Protocols {
		if err := p.Init(srv); err != nil {
			panic(err)
		}

		if err := p.Mount(srv.Router); err != nil {
			panic(err)
		}
	}
}
