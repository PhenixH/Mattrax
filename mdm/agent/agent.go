package agent

import (
	mattrax "github.com/mattrax/Mattrax/internal"
)

// Mount initialise the MDM server
func Mount(srv *mattrax.Server) {
	srv.Router.HandleFunc("/Manage/Sync.svc", Sync(srv)).Name("agent-sync").Methods("POST")
	srv.GlobalRouter.HandleFunc("/Manage/Enroll.svc", Enroll(srv)).Name("agent-enroll").Methods("GET", "POST")
}
