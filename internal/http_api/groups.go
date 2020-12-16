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
	"github.com/openzipkin/zipkin-go"
)

func Groups(srv *mattrax.Server) http.HandlerFunc {
	type CreateRequest struct {
		Name string `json:"name" validate:"required,alphanumspace,min=1,max=100"`
	}

	type CreateResponse struct {
		GroupID string `json:"group_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)
		if r.Method == http.MethodGet {
			limit, offset, err := middleware.GetPaginationParams(r.URL.Query())
			if err != nil {
				span.Tag("warn", fmt.Sprintf("%s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			span.Tag("limit", fmt.Sprintf("%v", limit))
			span.Tag("offset", fmt.Sprintf("%v", offset))

			groups, err := srv.DB.WithTx(tx).GetGroups(r.Context(), db.GetGroupsParams{
				TenantID: vars["tenant"],
				Limit:    limit,
				Offset:   offset,
			})
			if err != nil {
				log.Printf("[GetGroups Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving groups: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if groups == nil {
				groups = make([]db.GetGroupsRow, 0)
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(groups); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			var cmd CreateRequest
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err := validate.Struct(cmd); err != nil {
				span.Tag("err", fmt.Sprintf("error validing CreateGroupRequest: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if vars["tenant"] == "" {
				span.Tag("err", "no tenant was specified")
				w.WriteHeader(http.StatusNotFound)
				return
			}

			groupID, err := srv.DB.WithTx(tx).NewGroup(r.Context(), db.NewGroupParams{
				Name:     cmd.Name,
				TenantID: vars["tenant"],
			})
			if err != nil {
				if pqe, ok := err.(*pq.Error); ok && string(pqe.Code) == "23505" {
					span.Tag("warn", fmt.Sprintf("error creating new user due to unique constraint violation: %s", err))
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
				log.Printf("[CreateGroup Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error creating new group: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(CreateResponse{
				GroupID: groupID,
			}); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

func Group(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)
		if r.Method == http.MethodGet {
			user, err := srv.DB.WithTx(tx).GetGroup(r.Context(), db.GetGroupParams{
				ID:       vars["gid"],
				TenantID: vars["tenant"],
			})
			if err == sql.ErrNoRows {
				span.Tag("warn", "group not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[GetGroup Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("error retrieving group: %s", err))
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
			var cmd db.Group
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			query := `UPDATE groups SET name=COALESCE(NULLIF($3, ''), name) WHERE id = $1 AND tenant_id=$2;`
			if _, err := tx.Exec(query, vars["gid"], vars["tenant"], cmd.Name); err == sql.ErrNoRows {
				span.Tag("warn", "group not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[UpdateGroup Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error updating group: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == http.MethodDelete {
			err := srv.DB.WithTx(tx).DeleteGroup(r.Context(), db.DeleteGroupParams{
				ID:       vars["gid"],
				TenantID: vars["tenant"],
			})
			if err == sql.ErrNoRows {
				span.Tag("warn", "group not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[DeleteGroup Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("error deleting group: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func GroupPolicies(srv *mattrax.Server) http.HandlerFunc {
	type Request struct {
		Policies []string `json:"policies"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)
		if r.Method == http.MethodGet {
			devices, err := srv.DB.WithTx(tx).GetPoliciesInGroup(r.Context(), db.GetPoliciesInGroupParams{
				TenantID: vars["tenant"],
				GroupID:  vars["gid"],
				// TODO: Pagination
				Limit:  100,
				Offset: 0,
			})
			if err == sql.ErrNoRows {
				span.Tag("warn", "group not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[GetPoliciesInGroup Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("error retrieving policies in groups: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if devices == nil {
				devices = make([]db.GetPoliciesInGroupRow, 0)
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(devices); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			var cmd Request
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			for _, policyID := range cmd.Policies {
				fmt.Println(vars["gid"], policyID)
				if err := srv.DB.WithTx(tx).AddPolicyToGroup(r.Context(), db.AddPolicyToGroupParams{
					GroupID:  vars["gid"],
					PolicyID: policyID,
				}); err != nil {
					log.Printf("[AddPolicyToGroup Error]: %s\n", err)
					span.Tag("warn", fmt.Sprintf("error adding policy to group: %s", err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == http.MethodDelete {
			var cmd Request
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			for _, policyID := range cmd.Policies {
				if err := srv.DB.WithTx(tx).RemovePolicyFromGroup(r.Context(), db.RemovePolicyFromGroupParams{
					GroupID:  vars["gid"],
					PolicyID: policyID,
				}); err == sql.ErrNoRows {
					// TODO: Does this need to be handled
					span.Tag("warn", "policy group not found")
					w.WriteHeader(http.StatusNotFound)
					return
				} else if err != nil {
					log.Printf("[RemoveDeviceFromGroup Error]: %s\n", err)
					span.Tag("warn", fmt.Sprintf("error removing policy from group: %s", err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func GroupDevices(srv *mattrax.Server) http.HandlerFunc {
	type Request struct {
		Devices []string `json:"devices"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)
		if r.Method == http.MethodGet {
			devices, err := srv.DB.WithTx(tx).GetDevicesInGroup(r.Context(), db.GetDevicesInGroupParams{
				TenantID: vars["tenant"],
				GroupID:  vars["gid"],
				// TODO: Pagination
				Limit:  100,
				Offset: 0,
			})
			if err == sql.ErrNoRows {
				span.Tag("warn", "group not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[GetDevicesInGroup Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("error retrieving devices in groups: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if devices == nil {
				devices = make([]db.GetDevicesInGroupRow, 0)
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(devices); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodPost {
			var cmd Request
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			for _, deviceID := range cmd.Devices {
				if err := srv.DB.WithTx(tx).AddDeviceToGroup(r.Context(), db.AddDeviceToGroupParams{
					GroupID:  vars["gid"],
					DeviceID: deviceID,
				}); err != nil {
					log.Printf("[AddDevicesToGroup Error]: %s\n", err)
					span.Tag("warn", fmt.Sprintf("error adding device to group: %s", err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == http.MethodDelete {
			var cmd Request
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			for _, deviceID := range cmd.Devices {
				if err := srv.DB.WithTx(tx).RemoveDeviceFromGroup(r.Context(), db.RemoveDeviceFromGroupParams{
					GroupID:  vars["gid"],
					DeviceID: deviceID,
				}); err == sql.ErrNoRows {
					// TODO: Does this need to be handled
					span.Tag("warn", "device group not found")
					w.WriteHeader(http.StatusNotFound)
					return
				} else if err != nil {
					log.Printf("[RemoveDeviceFromGroup Error]: %s\n", err)
					span.Tag("warn", fmt.Sprintf("error removing device from group: %s", err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}
