package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		span := zipkin.SpanOrNoopFromContext(r.Context())
		claims := middleware.AuthClaimsFromContext(r.Context())
		if claims == nil {
			span.Tag("warn", "authentication claims are not set")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.Method == http.MethodGet {
			tenants, err := srv.DB.GetUserTenants(r.Context(), claims.Subject)
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

			tenantID, err := srv.DB.NewTenant(r.Context(), db.NewTenantParams(cmd))
			if err != nil {
				log.Printf("[NewTenant Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error creating new tenant: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := srv.DB.ScopeUserToTenant(r.Context(), db.ScopeUserToTenantParams{
				UserUpn:         claims.Subject,
				TenantID:        tenantID,
				PermissionLevel: db.UserPermissionLevelAdministrator,
			}); err != nil {
				log.Printf("[ScopeUserToTenant Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error scoping user to tenant: %s", err))
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
