package http_api

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/internal/middleware"
	"github.com/openzipkin/zipkin-go"
)

func Tenants(srv *mattrax.Server) http.HandlerFunc {
	type Request struct {
		DisplayName   string `json:"display_name" validate:"required,alphanumspace,min=1,max=100"`
		PrimaryDomain string `json:"primary_domain" validate:"required,fqdn,min=1,max=100"`
	}

	type Response struct {
		TenantID string `json:"tenant_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		claims := middleware.AuthClaimsFromContext(r.Context())
		if claims == nil {
			span.Tag("warn", "authentication claims are not set")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.Method == http.MethodGet {
			tenants, err := srv.DB.WithTx(tx).GetUserTenants(r.Context(), claims.Subject)
			if err != nil {
				log.Printf("[GetUserTenants Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving tenants: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(tenants); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			var cmd Request
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err := validate.Struct(cmd); err != nil {
				span.Tag("err", fmt.Sprintf("error validing new tenant request: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			tenantID, err := srv.DB.WithTx(tx).NewTenant(r.Context(), db.NewTenantParams(cmd))
			if err != nil {
				log.Printf("[NewTenant Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error creating new tenant: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := srv.DB.WithTx(tx).ScopeUserToTenant(r.Context(), db.ScopeUserToTenantParams{
				UserUpn:         claims.Subject,
				TenantID:        tenantID,
				PermissionLevel: db.UserPermissionLevelAdministrator,
			}); err != nil {
				log.Printf("[ScopeUserToTenant Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error scoping user to tenant: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if _, err := srv.DB.WithTx(tx).AddDomainToTenant(r.Context(), db.AddDomainToTenantParams{
				TenantID: tenantID,
				Domain:   cmd.PrimaryDomain,
			}); err != nil {
				log.Printf("[AddDomainToTenant Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error adding default domain to new tenant: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(Response{
				TenantID: tenantID,
			}); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func TenantDomain(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)

		if r.Method == http.MethodPost {
			linkingCode, err := srv.DB.WithTx(tx).AddDomainToTenant(r.Context(), db.AddDomainToTenantParams{
				Domain:   vars["domain"],
				TenantID: vars["tenant"],
			})
			if err != nil {
				log.Printf("[AddDomainToTenant Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error adding domain to tenant: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(db.GetTenantDomainsRow{
				Domain:      vars["domain"],
				LinkingCode: linkingCode,
				Verified:    false,
			}); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPatch {
			domain, err := srv.DB.WithTx(tx).GetTenantDomain(r.Context(), db.GetTenantDomainParams{
				Domain:   vars["domain"],
				TenantID: vars["tenant"],
			})
			if err != nil {
				log.Printf("[GetTenantDomain Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error getting tenant domain: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			records, err := net.LookupTXT(vars["domain"])
			if err != nil {
				log.Printf("[TXT Record Lookup Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error looking up txt record: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var verified = false
			for _, record := range records {
				if strings.HasPrefix(record, "mttx") && strings.TrimPrefix(record, "mttx") == domain.LinkingCode {
					verified = true
				}
			}

			if domain.Verified != verified {
				if err := srv.DB.WithTx(tx).UpdateDomain(r.Context(), db.UpdateDomainParams{
					Domain:   vars["domain"],
					TenantID: vars["tenant"],
					Verified: verified,
				}); err != nil {
					log.Printf("[UpdateDomain Error]: %s\n", err)
					span.Tag("err", fmt.Sprintf("error updating tenant domain: %s", err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(verified); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodDelete {
			if err := srv.DB.WithTx(tx).DeleteDomain(r.Context(), db.DeleteDomainParams{
				Domain:   vars["domain"],
				TenantID: vars["tenant"],
			}); err != nil {
				log.Printf("[DeleteDomain Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error deleting tenant domain: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}
