package windows

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/mattrax/Mattrax/pkg/syncml"
	"github.com/openzipkin/zipkin-go"
)

// Manage handles the continued management of the device.
func Manage(p *Protocol) http.HandlerFunc {
	var Handler = func(span zipkin.Span, w http.ResponseWriter, r *http.Request, cmd syncml.Message) syncml.Response {
		var res = syncml.NewResponse(cmd)
		span.Tag("device_id", cmd.Header.SourceURI)

		if url, err := url.ParseRequestURI(cmd.Header.TargetURI); err != nil || url.Host != p.srv.Args.Domain {
			err = fmt.Errorf("the request was not destined for this server")
			log.Printf("[Error] Manage | %s\n", err)
			span.Tag("err", fmt.Sprintf("%s", err))
			return syncml.NewBlankResponse(cmd, syncml.StatusBadRequest)
		}

		for _, command := range cmd.Body.Commands {
			if command.XMLName.Local == "Alert" && command.Data == "1226" && len(command.Body) == 1 {
				if item := command.Body[0]; item.Meta.Type == "com.microsoft:mdm.unenrollment.userrequest" && item.Meta.Format == "int" && item.Data == "1" {
					fmt.Println("User initiated unenroll triggered for device:", cmd.Header.SourceURI)
				}
			}
			// TODO: Update device information that was missing from enrollment
		}

		return res
	}

	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())

		var cmd syncml.Message
		if err := syncml.Read(&cmd, w, r); err != nil {
			err = fmt.Errorf("error marshaling syncml body from client: %w", err)
			log.Printf("[Error] Manage | %v\n", err)
			span.Tag("err", fmt.Sprintf("%s", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := Handler(span, w, r, cmd).Respond(w); err != nil {
			err = fmt.Errorf("error sending response to client: %w", err)
			log.Printf("[Error] Manage | %v\n", err)
			span.Tag("err", fmt.Sprintf("%s", err))
			return
		}
	}
}
