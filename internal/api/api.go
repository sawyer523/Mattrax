package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattrax/Mattrax/internal/datastore"
	"github.com/mattrax/Mattrax/internal/mattrax/settings"
	"github.com/mattrax/Mattrax/internal/services/device"
	"github.com/mattrax/Mattrax/internal/services/policy"
)

func Middleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			next.ServeHTTP(w, req)
		})
	}
}

func IndexHandler(version string) http.HandlerFunc {
	res := struct {
		Status         string `json:"status"`
		MattraxVersion string `json:"mattrax_version"`
		Docs           string `json:"docs"`
	}{
		Status:         "ok",
		MattraxVersion: version,
		Docs:           "https://github.com/mattrax/Mattrax/wiki/API",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(res); err != nil {
			panic(err) // TODO
		}
	}
}

type APIError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func SettingsHandler(ss *settings.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			var cmd settings.Settings
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				panic(err) // TODO
			}
			if err := ss.Set(cmd); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  err.Error(),
				}); err != nil {
					panic(err) // TODO
				}
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := json.NewEncoder(w).Encode(ss.Get()); err != nil {
			panic(err) // TODO
		}
	}
}

func DeviceHandler(ds *device.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if r.Method == http.MethodPatch {
			var cmd device.Device
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				panic(err) // TODO
			}
			if err := ds.Save(vars["uuid"], cmd); err == datastore.ErrNotFound {
				if err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  "device not found with specified uuid",
				}); err != nil {
					panic(err) // TODO
				}
			} else if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  err.Error(),
				})
				if err != nil {
					panic(err) // TODO
				}
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodDelete {
			if err := ds.Unenroll(vars["uuid"]); err == datastore.ErrNotFound {
				if err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  "device not found with specified uuid",
				}); err != nil {
					panic(err) // TODO
				}
			} else if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  err.Error(),
				})
				if err != nil {
					panic(err) // TODO
				}
				return
			}
			w.WriteHeader(http.StatusAccepted)
			return
		}

		if uuid, ok := vars["uuid"]; ok {
			device, err := ds.Get(uuid)
			if err == datastore.ErrNotFound {
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  "device not found",
				}); err != nil {
					panic(err) // TODO
				}
				return
			} else if err != nil {
				panic(err) // TODO
			}

			if err := json.NewEncoder(w).Encode(device); err != nil {
				panic(err) // TODO
			}
		} else {
			devices, err := ds.GetAll()
			if err != nil {
				panic(err) // TODO
			}

			if err := json.NewEncoder(w).Encode(devices); err != nil {
				panic(err) // TODO
			}
		}
	}
}

func PolicyHandler(ps *policy.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if r.Method == http.MethodPatch {
			var cmd policy.Policy
			if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
				panic(err) // TODO
			}
			if err := ps.Save(vars["uuid"], cmd); err == datastore.ErrNotFound {
				if err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  "policy not found with specified uuid",
				}); err != nil {
					panic(err) // TODO
				}
			} else if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  err.Error(),
				})
				if err != nil {
					panic(err) // TODO
				}
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.Method == http.MethodDelete {
			if err := ps.Delete(vars["uuid"]); err == datastore.ErrNotFound {
				if err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  "policy not found with specified uuid",
				}); err != nil {
					panic(err) // TODO
				}
			} else if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  err.Error(),
				})
				if err != nil {
					panic(err) // TODO
				}
				return
			}
			w.WriteHeader(http.StatusAccepted)
			return
		}

		if uuid, ok := vars["uuid"]; ok {
			device, err := ps.Get(uuid)
			if err == datastore.ErrNotFound {
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(APIError{
					Status: "error",
					Error:  "device not found",
				}); err != nil {
					panic(err) // TODO
				}
				return
			} else if err != nil {
				panic(err) // TODO
			}

			if err := json.NewEncoder(w).Encode(device); err != nil {
				panic(err) // TODO
			}
		} else {
			devices, err := ps.GetAll()
			if err != nil {
				panic(err) // TODO
			}

			if err := json.NewEncoder(w).Encode(devices); err != nil {
				panic(err) // TODO
			}
		}
	}
}
