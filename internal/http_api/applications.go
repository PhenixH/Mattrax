package http_api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/internal/middleware"
	"github.com/mattrax/Mattrax/pkg/null"
	"github.com/openzipkin/zipkin-go"
)

func Applications(srv *mattrax.Server) http.HandlerFunc {
	type CreateRequest struct {
		Name null.String `json:"name" validate:"required,alphanumspace,min=1,max=100"`
	}

	type CreateResponse struct {
		ApplicationID string `json:"application_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)
		if r.Method == http.MethodGet {
			limit, offset, err := middleware.GetPaginationParams(r.URL.Query())
			if err != nil {
				span.Tag("err", fmt.Sprintf("%s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			span.Tag("limit", fmt.Sprintf("%v", limit))
			span.Tag("offset", fmt.Sprintf("%v", offset))

			applications, err := srv.DB.WithTx(tx).GetApplications(r.Context(), db.GetApplicationsParams{
				TenantID: vars["tenant"],
				Limit:    limit,
				Offset:   offset,
			})
			if err != nil {
				log.Printf("[GetApplications Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving applications: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if applications == nil {
				applications = make([]db.GetApplicationsRow, 0)
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(applications); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			var cmd CreateRequest
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err := validate.Struct(cmd); err != nil {
				span.Tag("err", fmt.Sprintf("error validing CreateGroupRequest: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if vars["tenant"] == "" {
				span.Tag("err", fmt.Sprintf("no tenant was specified"))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			applicationID, err := srv.DB.WithTx(tx).NewApplication(r.Context(), db.NewApplicationParams{
				Name:     cmd.Name,
				TenantID: vars["tenant"],
			})
			if err != nil {
				if pqe, ok := err.(*pq.Error); ok && string(pqe.Code) == "23505" {
					span.Tag("warn", fmt.Sprintf("error creating new user due to unique constraint violation: %s", err))
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}

				log.Printf("[NewApplication Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error creating new application: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(CreateResponse{
				ApplicationID: applicationID,
			}); err != nil {
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func Application(srv *mattrax.Server) http.HandlerFunc {
	type ApplicationResponse struct {
		db.GetApplicationRow
		Targets []db.GetApplicationTargetsRow `json:"targets"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)

		if r.Method == http.MethodGet {
			app, err := srv.DB.WithTx(tx).GetApplication(r.Context(), db.GetApplicationParams{
				ID:       vars["aid"],
				TenantID: vars["tenant"],
			})
			if err == sql.ErrNoRows {
				span.Tag("warn", "application not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[GetApplication Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving application: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			appTargets, err := srv.DB.WithTx(tx).GetApplicationTargets(r.Context(), db.GetApplicationTargetsParams{
				AppID:    vars["aid"],
				TenantID: vars["tenant"],
			})
			if err != nil {
				log.Printf("[GetApplicationTargets Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving application targets: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(ApplicationResponse{app, appTargets}); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPatch {
			var cmd db.Application
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			err := srv.DB.WithTx(tx).UpdateApplication(r.Context(), db.UpdateApplicationParams{
				ID:          vars["aid"],
				TenantID:    vars["tenant"],
				Name:        cmd.Name,
				Description: cmd.Description,
				Publisher:   cmd.Publisher,
			})
			if err == sql.ErrNoRows {
				span.Tag("warn", "application not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[UpdateApplication Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error updating application: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == http.MethodDelete {
			err := srv.DB.WithTx(tx).DeleteApplication(r.Context(), db.DeleteApplicationParams{
				ID:       vars["aid"],
				TenantID: vars["tenant"],
			})
			if err == sql.ErrNoRows {
				span.Tag("warn", "application not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[DeleteApplication Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("error deleting application: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}
