package api

import (
	"net/http"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
)

const MaxJSONBodySize = 2097152

// Mount initialises the API
func Mount(srv *mattrax.Server) {
	r := srv.Router.PathPrefix("/api").Subrouter()
	r.Use(Headers(srv))
	r.Use(mux.CORSMethodMiddleware(r))

	r.HandleFunc("/login", Login(srv)).Methods(http.MethodPost, http.MethodOptions)

	rAuthed := r.PathPrefix("/").Subrouter()
	rAuthed.Use(RequireAuthentication(srv))

	rAuthed.HandleFunc("/devices", Devices(srv)).Methods(http.MethodGet, http.MethodOptions).Name("/devices")
	rAuthed.HandleFunc("/device/{id}", Device(srv)).Methods(http.MethodGet, http.MethodOptions).Name("/devices/:id")
	rAuthed.HandleFunc("/device/{id}/info", DeviceInformation(srv)).Methods(http.MethodGet, http.MethodOptions).Name("/devices/:id/info")
	rAuthed.HandleFunc("/device/{id}/scope", DeviceScope(srv)).Methods(http.MethodGet, http.MethodOptions).Name("/devices/:id/scope")
	rAuthed.HandleFunc("/groups", Groups(srv)).Methods(http.MethodGet, http.MethodOptions).Name("/groups")
	rAuthed.HandleFunc("/group/{id}", Group(srv)).Methods(http.MethodGet, http.MethodOptions).Name("/groups/:id")
	rAuthed.HandleFunc("/policies", Policies(srv)).Methods(http.MethodGet, http.MethodOptions).Name("/policies")
	rAuthed.HandleFunc("/policy/{id}", Policy(srv)).Methods(http.MethodGet, http.MethodOptions).Name("/policies/:id")
	rAuthed.HandleFunc("/users", Users(srv)).Methods(http.MethodGet, http.MethodPost, http.MethodOptions).Name("/users")
	rAuthed.HandleFunc("/user/{upn}", User(srv)).Methods(http.MethodGet, http.MethodOptions).Name("/users/:upn")
}
