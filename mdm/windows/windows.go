package windows

import (
	mattrax "github.com/mattrax/Mattrax/internal"
)

// EnrollmentPolicyID is the unique ID of the MDM server's enrollment policy
const EnrollmentPolicyID = "mattrax-identity"

// EnrollmentPolicyFriendlyName is the friendly name of the server's enrollment policy
const EnrollmentPolicyFriendlyName = "Mattrax Identity Certificate Policy"

// ProviderID is the unique ID used to identify the MDM server to the management client
// This is not shown to the user and should be globally identifying. It is required by some CSP's.
const ProviderID = "Mattrax"

// Protocol contains the handlers for the Base protocol
type Protocol struct {
	srv *mattrax.Server
}

// TODO: Unused remove
func (p *Protocol) ID() string {
	return "windows"
}

func (p *Protocol) Init(srv *mattrax.Server) (err error) {
	p.srv = srv
	return nil
}

type status struct{}

func (p *Protocol) Status() (interface{}, error) {
	return status{}, nil
}
