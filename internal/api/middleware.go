package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/openzipkin/zipkin-go"
)

func Headers(srv *mattrax.Server) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if srv.Args.Debug {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			// w.Header().Set("Strict-Transport-Security", "max-age=31536000;")
			w.Header().Set("Content-Security-Policy", "")
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

func RequireAuthentication(srv *mattrax.Server) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := zipkin.SpanOrNoopFromContext(r.Context())

			authorizationHeader := strings.Split(r.Header.Get("Authorization"), " ")
			if len(authorizationHeader) != 2 {
				span.Tag("error", "Authorization header invalid format")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

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

			span.Tag("upn", claims.Subject)
			span.Tag("org", claims.Organisation)
			next.ServeHTTP(w, r)
		})
	}
}
