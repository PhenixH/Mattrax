package syncml

import (
	"fmt"
	"net/http"

	"github.com/mattrax/xml"
)

// MaxRequestBodySize is the maximum amount of data that is allowed in a single request
const MaxRequestBodySize = 524288

// Read safely decodes a SyncML request from the HTTP body into a struct
func Read(v *Message, w http.ResponseWriter, r *http.Request) error {
	if r.ContentLength > MaxRequestBodySize {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return fmt.Errorf("Error decoding request of type '%T': Request body of size '%d' is larger than the maximum supported size of '%d'", v, r.ContentLength, MaxRequestBodySize)
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)
	if err := xml.NewDecoder(r.Body).Decode(v); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("error decoding request of type '%T': %v", v, err)
	}

	return nil
}
