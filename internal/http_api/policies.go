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
	"github.com/mattrax/Mattrax/mdm"
	"github.com/openzipkin/zipkin-go"
)

func Policies(srv *mattrax.Server) http.HandlerFunc {
	type CreateRequest struct {
		Name string `json:"name" validate:"required,alphanumspace,min=1,max=100"`
		Type string `json:"type"` // TODO: Validate against JSON
	}

	type CreateResponse struct {
		PolicyID string `json:"policy_id"`
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

			policies, err := srv.DB.WithTx(tx).GetPolicies(r.Context(), db.GetPoliciesParams{
				TenantID: vars["tenant"],
				Limit:    limit,
				Offset:   offset,
			})
			if err != nil {
				log.Printf("[GetPolicies Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving policies: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if policies == nil {
				policies = make([]db.GetPoliciesRow, 0)
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(policies); err != nil {
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
			policyID, err := srv.DB.WithTx(tx).NewPolicy(r.Context(), db.NewPolicyParams{
				Name:     cmd.Name,
				TenantID: vars["tenant"],
				Type:     cmd.Type,
			})
			if err != nil {
				if pqe, ok := err.(*pq.Error); ok && string(pqe.Code) == "23505" {
					span.Tag("warn", fmt.Sprintf("error creating new user due to unique constraint violation: %s", err))
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}

				log.Printf("[NewPolicy Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error creating new policy: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(CreateResponse{
				PolicyID: policyID,
			}); err != nil {
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func Policy(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)

		if r.Method == http.MethodGet {
			policy, err := srv.DB.WithTx(tx).GetPolicy(r.Context(), db.GetPolicyParams{
				ID:       vars["pid"],
				TenantID: vars["tenant"],
			})
			if err == sql.ErrNoRows {
				span.Tag("warn", "policy not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[GetPolicy Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving policy: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(policy); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPatch {
			var cmd db.Policy
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if len(cmd.Payload) == 0 {
				cmd.Payload = []byte("{}")
			}

			query := `UPDATE policies SET name=COALESCE(NULLIF($3, ''), name), type=COALESCE(NULLIF($4, ''), type), payload=payload||$5 WHERE id = $1 AND tenant_id=$2;`
			if _, err := tx.Exec(query, vars["pid"], vars["tenant"], cmd.Name, cmd.Type, string(cmd.Payload)); err == sql.ErrNoRows {
				span.Tag("warn", "policy not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[UpdatePolicy Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error updating policy: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// TEMP
			policyRow, err := srv.DB.WithTx(tx).GetPolicy(r.Context(), db.GetPolicyParams{
				ID:       vars["pid"],
				TenantID: vars["tenant"],
			})
			if err != nil {
				log.Printf("[GetPolicy Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving policy: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var policy = db.Policy{
				ID:          vars["pid"],
				TenantID:    vars["tenant"],
				Name:        policyRow.Name,
				Type:        policyRow.Type,
				Payload:     policyRow.Payload,
				Description: policyRow.Description,
			}

			for _, p := range mdm.Protocols {
				events := p.Events()
				if events.UpdatePolicy != nil {
					err := events.UpdatePolicy(policy)
					if err != nil {
						panic(err) // TODO
					}
				}
			}
			// END TEMP

			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == http.MethodDelete {
			err := srv.DB.WithTx(tx).DeletePolicy(r.Context(), db.DeletePolicyParams{
				ID:       vars["pid"],
				TenantID: vars["tenant"],
			})
			if err == sql.ErrNoRows {
				span.Tag("warn", "policy not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[DeletePolicy Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("error deleting policy: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func PolicyScope(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)

		groups, err := srv.DB.WithTx(tx).GetPolicyGroups(r.Context(), vars["pid"])
		if err == sql.ErrNoRows {
			span.Tag("warn", "policy groups not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetPolicyGroups Error]: %s\n", err)
			span.Tag("err", fmt.Sprintf("error retrieving policy groups: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if groups == nil {
			groups = make([]db.GetPolicyGroupsRow, 0)
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(groups); err != nil {
			span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
