package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/authentication"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/openzipkin/zipkin-go"
	"golang.org/x/crypto/bcrypt"
)

func Login(srv *mattrax.Server) http.HandlerFunc {
	type Request struct {
		UPN      string `json:"upn"`
		Password string `json:"password"`
	}

	type Response struct {
		Token   string                 `json:"token"`
		Tenants []db.GetUserTenantsRow `json:"tenants,notnull"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())

		var cmd Request
		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		span.Tag("upn", cmd.UPN)

		var audience = "dashboard" // TODO: Set to enrollment if device enrolling process

		user, err := srv.DB.GetUserSecure(r.Context(), cmd.UPN)
		if err != nil {
			log.Printf("[GetUserSecure Error]: %s\n", err)
			span.Tag("err", fmt.Sprintf("[GetUserSecure Error]: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !user.Password.Valid {
			span.Tag("err", fmt.Sprintf("user does not have a password set"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if audience != "enrollment" && user.TenantID.Valid == true {
			span.Tag("err", fmt.Sprintf("user does not have tenant permission to login to dashboard"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(cmd.Password)); err == bcrypt.ErrMismatchedHashAndPassword {
			span.Tag("warn", fmt.Sprintf("user password did not match hash"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else if err != nil {
			span.Tag("err", fmt.Sprintf("error comparing password to hash: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authToken, _, err := srv.Auth.IssueToken(audience, authentication.AuthClaims{
			Subject:  cmd.UPN,
			FullName: user.Fullname,
		})
		if err != nil {
			log.Printf("[IssueToken Error]: %s\n", err)
			span.Tag("err", fmt.Sprintf("error issuing auth token: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tenants, err := srv.DB.GetUserTenants(r.Context(), cmd.UPN)
		if err != nil {
			log.Printf("[GetUserTenants Error]: %s\n", err)
			span.Tag("err", fmt.Sprintf("error retrieving users tenants: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if tenants == nil {
			tenants = make([]db.GetUserTenantsRow, 0)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(Response{
			Token:   authToken,
			Tenants: tenants,
		}); err != nil {
			span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
