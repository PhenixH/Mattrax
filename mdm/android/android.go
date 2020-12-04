package android

import (
	"context"
	"fmt"

	mattrax "github.com/mattrax/Mattrax/internal"
	"google.golang.org/api/androidmanagement/v1"
	"google.golang.org/api/option"
)

// Protocol contains the handlers for the Base protocol
type Protocol struct {
	srv *mattrax.Server
	ams *androidmanagement.Service
}

// TODO: Unused remove
func (p *Protocol) ID() string {
	return "base"
}

func (p *Protocol) Init(srv *mattrax.Server) (err error) {
	if srv.Args.GoogleServiceAccountPath == "" {
		return fmt.Errorf("error starting Android management service: no service account configured")
	}

	if p.ams, err = androidmanagement.NewService(context.Background(), option.WithCredentialsFile(srv.Args.GoogleServiceAccountPath)); err != nil {
		return fmt.Errorf("error creating android management API service: %w", err)
	}
	p.srv = srv

	return nil
}

func (p *Protocol) Mount() error {
	return nil
}
