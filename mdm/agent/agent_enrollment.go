package agent

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/pkg"
	"github.com/rs/zerolog/log"
)

// EnrollRequest is the request for a device to enroll
type EnrollRequest struct {
	UDID               string `json:"udid"`
	Hostname           string `json:"hostname"`
	CertificateRequest []byte `json:"certificateRequest"`
}

// EnrollResponse is the response with the device enrollment information
type EnrollResponse struct {
	Server          string `json:"server"`
	RootCertificate string `json:"rootCertificate"`
	Certificate     string `json:"certificate"`
}

// Enroll handles the device enrollment.
func Enroll(srv *mattrax.Server) http.HandlerFunc {
	syncServiceURL, err := pkg.GetNamedRouteURL(srv.GlobalRouter, "agent-sync")
	if err != nil {
		log.Fatal().Err(err).Msg("Error determining route URL") // TODO: Move error handling to main package
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if strings.Split(r.Header.Get("User-Agent"), " ")[0] != "MattraxAgent" || r.Header.Get("Content-Type") != "application/json+dm" {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		// TODO: Header based authentication

		var cmd EnrollRequest
		err := json.NewDecoder(r.Body).Decode(&cmd)
		if err != nil {
			log.Error().Err(err).Str("protocol", "agent").Msg("Error decoding enrollment request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		csr, err := x509.ParseCertificateRequest(cmd.CertificateRequest)
		if err != nil {
			log.Error().Err(err).Str("protocol", "agent").Msg("Error decoding certificate signing request")
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err = csr.CheckSignature(); err != nil {
			log.Error().Err(err).Str("protocol", "agent").Msg("Invalid certificate signing request signature")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		identityCertificate, _, signedClientCertificate, err := srv.Cert.IdentitySignCSR(csr, pkix.Name{
			CommonName:   cmd.UDID,
			Organization: []string{"Mattrax Agent"},
		})
		if err != nil {
			log.Error().Err(err).Str("protocol", "agent").Msg("error creating client certificate")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var res = EnrollResponse{
			Server:          syncServiceURL,
			RootCertificate: base64.StdEncoding.EncodeToString(identityCertificate.Raw),
			Certificate:     base64.StdEncoding.EncodeToString(signedClientCertificate),
		}

		w.Header().Set("Content-Type", "application/json+dm")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Error().Err(err).Str("protocol", "agent").Msg("error encoding & writting json response body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
