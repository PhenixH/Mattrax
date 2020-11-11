package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/db"
	"github.com/mattrax/Mattrax/internal/middleware"
	"github.com/openzipkin/zipkin-go"
)

func Devices(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)

		limit, offset, err := middleware.GetPaginationParams(r.URL.Query())
		if err != nil {
			span.Tag("warn", fmt.Sprintf("%s", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		span.Tag("limit", fmt.Sprintf("%v", limit))
		span.Tag("offset", fmt.Sprintf("%v", offset))

		devices, err := srv.DB.GetDevices(r.Context(), db.GetDevicesParams{
			TenantID: vars["tenant"],
			Limit:    limit,
			Offset:   offset,
		})
		if err != nil {
			log.Printf("[GetDevices Error]: %s\n", err)
			span.Tag("err", fmt.Sprintf("error getting devices: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(devices); err != nil {
			span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func Device(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)

		device, err := srv.DB.GetDevice(r.Context(), db.GetDeviceParams{
			ID:       vars["udid"],
			TenantID: vars["tenant"],
		})
		if err == sql.ErrNoRows {
			span.Tag("warn", "device not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetBasicDevice Error]: %s\n", err)
			span.Tag("err", fmt.Sprintf("error retrieving device: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(device); err != nil {
			span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func DeviceInformation(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)

		device, err := srv.DB.GetDevice(r.Context(), db.GetDeviceParams{
			ID:       vars["udid"],
			TenantID: vars["tenant"],
		})
		if err == sql.ErrNoRows {
			span.Tag("warn", "device not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetDevice Error]: %s\n", err)
			span.Tag("err", fmt.Sprintf("error retrieving device: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]map[string]interface{}{
			"Device Information": {
				"Computer Name": device.Name,
				// "Serial Number":
			},
			"Software Information": {
				// "Operating System":         "Windows 10", // TODO
				// "Operating System Version": device.OperatingSystem,
			},
			"MDM": {
				// "Last Seen":        device.Lastseen,
				// "Last Seen Status": device.LastseenStatus,
			},
		}); err != nil {
			span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func DeviceScope(srv *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanOrNoopFromContext(r.Context())
		vars := mux.Vars(r)

		groups, err := srv.DB.GetDeviceGroups(r.Context(), vars["udid"])
		if err == sql.ErrNoRows {
			span.Tag("warn", "device groups not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetBasicDeviceScopedGroups Error]: %s\n", err)
			span.Tag("err", fmt.Sprintf("error retrieving device groups: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		policies, err := srv.DB.GetDevicePolicies(r.Context(), vars["udid"])
		if err == sql.ErrNoRows {
			span.Tag("warn", "device policies not found")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("[GetBasicDeviceScopedPolicies Error]: %s\n", err)
			span.Tag("err", fmt.Sprintf("error retrieving device policies: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"groups":   groups,
			"policies": policies,
		}); err != nil {
			span.Tag("warn", fmt.Sprintf("error encoding JSON response: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
