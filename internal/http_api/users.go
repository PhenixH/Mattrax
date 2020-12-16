package http_api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/internal/middleware"
	"github.com/mattrax/Mattrax/pkg/null"
	"github.com/openzipkin/zipkin-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
)

func Users(srv *mattrax.Server) http.HandlerFunc {
	type CreateRequest struct {
		UPN      string `json:"upn" validate:"required,email,min=1,max=100"`
		FullName string `json:"fullname" validate:"required,alphanumspace,min=1,max=100"`
		Password string `json:"password" validate:"required,min=1,max=100"`
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

			users, err := srv.DB.WithTx(tx).GetUsersInTenant(r.Context(), db.GetUsersInTenantParams{
				TenantID: null.String{
					String: vars["tenant"],
					Valid:  true,
				},
				Limit:  limit,
				Offset: offset,
			})
			if err != nil {
				log.Printf("[GetUsers Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error retrieving users: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(users); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		} else if r.Method == http.MethodPost {
			var cmd CreateRequest
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if err := validate.Struct(cmd); err != nil {
				span.Tag("err", fmt.Sprintf("error validing CreateUserRequest: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			passwordHash, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), 15)
			if err != nil {
				log.Printf("[BcryptGenerateFromPassword Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error generating bcrypt hash for password: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err := srv.DB.WithTx(tx).NewUser(r.Context(), db.NewUserParams{
				UPN:      cmd.UPN,
				Fullname: cmd.FullName,
				Password: null.String{
					String: string(passwordHash),
					Valid:  true,
				},
				TenantID: null.String{
					String: vars["tenant"],
					Valid:  vars["tenant"] != "",
				},
			}); err != nil {
				if pqe, ok := err.(*pq.Error); ok && string(pqe.Code) == "23505" {
					span.Tag("warn", fmt.Sprintf("error creating new user due to unique constraint violation: %s", err))
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}

				log.Printf("[CreateUser Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error creating new user: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func User(srv *mattrax.Server) http.HandlerFunc {
	// TODO: Join with db.User
	type PatchRequest struct {
		UPN      string         `json:"upn"`
		Fullname string         `json:"fullname"`
		Disabled *bool          `json:"disabled"` // TODO: SQLC type override for this and then use db.User here
		Password sql.NullString `json:"password"`
	}

	type PatchResponse struct {
		UPN string `json:"upn"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tx := middleware.DBTxFromContext(r.Context())
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)

		if r.Method == http.MethodGet {
			user, err := srv.DB.WithTx(tx).GetUser(r.Context(), vars["upn"])
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
			var cmd PatchRequest
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				log.Printf("[JsonDecode Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("JSON decode error: %s", err))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// TODO: Work for both types of user -> admin or within tenant + Can modify some about themselves without admin

			var upn string
			query := `UPDATE users SET upn=COALESCE(NULLIF($2, ''), upn), fullname=COALESCE(NULLIF($3, ''), fullname), password=COALESCE(NULLIF($4, ''), password), disabled=COALESCE($5, disabled) WHERE upn = $1 RETURNING upn;`
			err := tx.QueryRow(query, vars["upn"], cmd.UPN, cmd.Fullname, cmd.Password, cmd.Disabled).Scan(&upn)
			if err == sql.ErrNoRows {
				span.Tag("warn", "user not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[UpdateUser Error]: %s\n", err)
				span.Tag("err", fmt.Sprintf("error updating user: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			if err := json.NewEncoder(w).Encode(PatchResponse{
				UPN: upn,
			}); err != nil {
				span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if r.Method == http.MethodDelete {
			fmt.Println(vars["upn"], vars["tenant"])
			err := srv.DB.WithTx(tx).DeleteUserInTenant(r.Context(), db.DeleteUserInTenantParams{
				UPN: vars["upn"],
				TenantID: null.String{
					Valid:  true,
					String: vars["tenant"],
				},
			})
			if err == sql.ErrNoRows {
				span.Tag("warn", "user not found")
				w.WriteHeader(http.StatusNotFound)
				return
			} else if err != nil {
				log.Printf("[DeleteUser Error]: %s\n", err)
				span.Tag("warn", fmt.Sprintf("error deleting user: %s", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}
