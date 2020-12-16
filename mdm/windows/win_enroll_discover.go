package windows

import (
	"fmt"
	"log"
	"net/http"
	// "net/url"
	// "os"

	"github.com/mattrax/Mattrax/pkg"
	"github.com/mattrax/Mattrax/pkg/soap"
	"github.com/openzipkin/zipkin-go"
)

// Discover handles the discovery phase of enrollment.
func Discover(p *Protocol) http.HandlerFunc {
	enrollmentPolicyServiceURL, err := pkg.GetNamedRouteURL(p.srv.Router, "policy")
	enrollmentServiceURL, err2 := pkg.GetNamedRouteURL(p.srv.Router, "enrollment")
	if err != nil || err2 != nil {
		log.Fatal("Error: Policy or Enrollment endpoint not found on the router.")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())

		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			return
		}

		var cmd soap.DiscoverRequest
		if err := soap.Read(&cmd, w, r); err != nil {
			log.Printf("[Error] Enrollment Discover | %v\n", err)
			span.Tag("err", fmt.Sprintf("%s", err))
			return
		}

		if cmd.Body.RequestVersion == "" {
			cmd.Body.RequestVersion = "4.0"
		}

		var res = soap.ResponseEnvelope{
			Header: soap.ResponseHeader{
				RelatesTo: cmd.Header.MessageID,
			},
			Body: soap.ResponseEnvelopeBody{
				Body: soap.DiscoverResponse{
					AuthPolicy:                 "OnPremise",
					EnrollmentVersion:          cmd.Body.RequestVersion,
					EnrollmentPolicyServiceURL: enrollmentPolicyServiceURL,
					EnrollmentServiceURL:       enrollmentServiceURL,
				},
			},
		}
		res.Populate("http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/DiscoverResponse")
		if err := soap.Respond(res, w); err != nil {
			log.Printf("[Error] Enrollment Discover | %v\n", err)
			span.Tag("err", fmt.Sprintf("%s", err))
			return
		}
	}
}
