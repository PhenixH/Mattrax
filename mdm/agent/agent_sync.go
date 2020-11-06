package agent

import (
	"net/http"

	mattrax "github.com/mattrax/Mattrax/internal"
)

// Sync handles the device management sync.
func Sync(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
}
