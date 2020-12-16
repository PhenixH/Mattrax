package windows

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/mattrax/Mattrax/pkg/soap"
	"github.com/openzipkin/zipkin-go"
)

// Policy instructs the client how the generate the identity certificate.
// This endpoint is part of the spec MS-XCEP.
func Policy(p *Protocol) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())

		var cmd soap.PolicyRequest
		if err := soap.Read(&cmd, w, r); err != nil {
			log.Printf("[Error] Enrollment Policy | %v\n", err)
			span.Tag("err", fmt.Sprintf("%s", err))
			return
		}

		var res soap.ResponseEnvelope
		if cmd.Header.Action != "http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy/IPolicy/GetPolicies" {
			err := fmt.Errorf("incorrect request action")
			log.Printf("[Error] Enrollment Enroll | %s\n", err)
			span.Tag("err", fmt.Sprintf("%s", err))
			res = soap.NewFault("s:Receiver", "s:MessageFormat", "InvalidEnrollmentData", "The request was not destined for this endpoint", span.Context().TraceID.String())
		} else if url, err := url.ParseRequestURI(cmd.Header.To); err != nil || url.Host != p.srv.Args.Domain {
			err = fmt.Errorf("the request was not destined for this server")
			log.Printf("[Error] Enrollment Enroll | %s\n", err)
			span.Tag("err", fmt.Sprintf("%s", err))
			res = soap.NewFault("s:Receiver", "s:MessageFormat", "InvalidEnrollmentData", "The request was not destined for this server", span.Context().TraceID.String())
		} else {
			res = soap.NewPolicyResponse(cmd.Header.MessageID, EnrollmentPolicyID, EnrollmentPolicyFriendlyName)
		}

		if err := soap.Respond(res, w); err != nil {
			log.Printf("[Error] Enrollment Policy | %v\n", err)
			span.Tag("err", fmt.Sprintf("%s", err))
			return
		}
	}
}
