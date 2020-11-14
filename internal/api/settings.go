package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/internal/middleware"
	"github.com/openzipkin/zipkin-go"
)

func SettingsTenant(srv *mattrax.Server) http.HandlerFunc {
	// TODO: Replace with db.Tenant once sql.NullString fixed
	type PatchRequest struct {
		ID            string `json:"id"`
		DisplayName   string `json:"display_name"`
		PrimaryDomain string `json:"primary_domain"`
		Email         string `json:"email"`
		Phone         string `json:"phone"`
		Description   string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)
		if r.Method == http.MethodGet {
			tenant, err := srv.DB.GetTenant(r.Context(), vars["tenant"])
			if err == sql.ErrNoRows {
				span.Tag("warn", "tenant not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[GetTenant Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving tenant: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(tenant); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPatch {
			var cmd PatchRequest
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			query := `UPDATE tenants SET display_name=COALESCE(NULLIF($2, ''), display_name), email=COALESCE(NULLIF($3, ''), email), phone=COALESCE(NULLIF($4, ''), phone) WHERE id = $1;`
			if _, err := srv.DBConn.Exec(query, vars["tenant"], cmd.DisplayName, cmd.Email, cmd.Phone); err == sql.ErrNoRows {
				span.Tag("warn", "tenant not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[UpdateTenant Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error updating tenant: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func SettingsMe(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())
		claims := middleware.AuthClaimsFromContext(r.Context())
		if claims == nil {
			span.Tag("warn", "authentication claims are not set")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if r.Method == http.MethodGet {
			user, err := srv.DB.GetUser(r.Context(), claims.Subject)
			if err == sql.ErrNoRows {
				span.Tag("warn", "user not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[GetUser Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving user: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(user); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPatch {
			var cmd db.User
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			query := `UPDATE users SET fullname=COALESCE(NULLIF($2, ''), fullname) WHERE upn = $1;`
			if _, err := srv.DBConn.Exec(query, claims.Subject, cmd.Fullname); err == sql.ErrNoRows {
				span.Tag("warn", "user not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[UpdateTenant Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error updating user: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}
