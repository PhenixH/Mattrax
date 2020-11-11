package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/internal/middleware"
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

			users, err := srv.DB.GetUsersInTenant(r.Context(), db.GetUsersInTenantParams{
				TenantID: sql.NullString{
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

			if err := srv.DB.NewUser(r.Context(), db.NewUserParams{
				Upn:      cmd.UPN,
				Fullname: cmd.FullName,
				Password: sql.NullString{
					String: string(passwordHash),
					Valid:  true,
				},
				TenantID: sql.NullString{
					String: vars["tenant"],
					Valid:  vars["tenant"] != "",
				},
			}); err != nil {
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
	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)
		user, err := srv.DB.GetUser(r.Context(), vars["upn"])
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
	}
}
