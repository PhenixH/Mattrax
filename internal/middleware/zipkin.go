package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/openzipkin/zipkin-go"
)

// ZipkinExtended extends Zipkin intergration by naming the spans the Gorilla Mux route name & setting the X-Request-ID header
func ZipkinExtended(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanFromContext(r.Context())
		if route := mux.CurrentRoute(r); route != nil {
			name := route.GetName()
			if name != "" {
				span.SetName(name)
			}
		}

		w.Header().Set("X-Request-ID", span.Context().TraceID.String())
		span.Tag("host", r.Host)

		next.ServeHTTP(w, r)
	})
}
