package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/authentication"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/openzipkin/zipkin-go"
)

const MaxJSONBodySize = 2097152

func APIHeaders(srv *mattrax.Server) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if srv.Args.Debug {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization,Access-Control-Request-Headers")
				w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE") // TEMP: Bypass for another bug
			}
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, MaxJSONBodySize)
			next.ServeHTTP(w, r)
		})
	}
}

type contextKey string

var (
	claimsContextKey = contextKey("auth-claims")
	dbTxContextKey   = contextKey("db-tx")
)

func AuthClaimsFromContext(ctx context.Context) *authentication.AuthClaims {
	v := ctx.Value(claimsContextKey)
	if v == nil {
		return nil
	}
	claims := v.(authentication.AuthClaims)
	return &claims
}

func DBTxFromContext(ctx context.Context) *sql.Tx {
	v := ctx.Value(dbTxContextKey)
	if v == nil {
		return nil
	}
	return v.(*sql.Tx)
}

func RequireAuthentication(srv *mattrax.Server) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := zipkin.SpanOrNoopFromContext(r.Context())
			vars := mux.Vars(r)

			tx, err := srv.DBConn.Begin()
			if err != nil {
				span.Tag("error", fmt.Sprintf("Error creating DB transaction: %s\n", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			authorizationHeader := strings.Split(r.Header.Get("Authorization"), " ")
			if len(authorizationHeader) != 2 {
				span.Tag("error", "Authorization header invalid format")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if os.Getenv("MATTRAX_CLOUD_AUTH") != "" && os.Getenv("MATTRAX_CLOUD_AUTH") == authorizationHeader[1] {
				if os.Getenv("MATTRAX_CLOUD_AUTH_IP") != "" && !strings.HasPrefix(r.RemoteAddr, os.Getenv("MATTRAX_CLOUD_AUTH_IP")) {
					span.Tag("err", "Mattrax Cloud token used from invalid IP. Token may have leaked!")
					return
				}
				span.Tag("user", "MATTRAX_CLOUD_AUTH")
			} else {
				claims, err := srv.Auth.Token(authorizationHeader[1])
				if err != nil {
					log.Printf("[Authentication Error]: %s\n", err)
					span.Tag("error", fmt.Sprintf("Error parsing authentication claims: %s\n", err))
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if claims.BasicClaims.Audience != "dashboard" {
					span.Tag("error", "Authentication claims are not authorized to access the dashboard")
					w.WriteHeader(http.StatusForbidden)
					return
				}

				r = r.WithContext(context.WithValue(r.Context(), claimsContextKey, claims))
				span.Tag("upn", claims.Subject)

				if tid, ok := vars["tenant"]; ok {
					permissionLevel, err := srv.DB.GetUserPermissionLevelForTenant(r.Context(), db.GetUserPermissionLevelForTenantParams{
						UserUpn:  claims.Subject,
						TenantID: tid,
					})
					if err == sql.ErrNoRows {
						span.Tag("err", fmt.Sprintf("user is not authorized to this tenant"))
						w.WriteHeader(http.StatusForbidden)
						return
					} else if err != nil {
						log.Printf("[GetUserPermissionLevelForTenant Error]: %s\n", err)
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					if permissionLevel != db.UserPermissionLevelAdministrator {
						span.Tag("err", fmt.Sprintf("user lacking permission for this resource"))
						w.WriteHeader(http.StatusForbidden)
						return
					}
				}

				// INSECURE: POTENTIAL SQL INJECTION VECTOR
				if _, err := tx.Exec("SET mattrax.upn = '" + claims.Subject + "';"); err != nil {
					span.Tag("error", fmt.Sprintf("Error setting mattrax upn on DB transaction: %s\n", err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			r = r.WithContext(context.WithValue(r.Context(), dbTxContextKey, tx))

			w.Header().Add("X-Retrieved-At", time.Now().UTC().Format(http.TimeFormat))

			next.ServeHTTP(w, r)

			if err := tx.Commit(); err != nil {
				span.Tag("error", fmt.Sprintf("Error committing DB transaction: %s\n", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		})
	}
}

func GetPaginationParams(v url.Values) (int32, int32, error) {
	var limit int32 = 25
	if slimit := v.Get("limit"); slimit != "" {
		ilimit, err := strconv.Atoi(slimit)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid limit parameter. Must be number: %w", err)
		}

		if ilimit <= 0 {
			limit = 1
		} else if ilimit > 50 {
			limit = 50
		} else {
			limit = int32(ilimit)
		}
	}

	var offset int32 = 0
	if soffset := v.Get("offset"); soffset != "" {
		ioffset, err := strconv.Atoi(soffset)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid offset parameter. Must be number: %w", err)
		}

		if ioffset <= 0 {
			offset = 1
		} else {
			offset = int32(ioffset)
		}
	}

	return limit, offset, nil
}
