package api

import (
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/settings"
	"gopkg.in/yaml.v2"
)

// Response is what is a generic HTTP response
type Response struct {
	Success bool
	Msg     string `json:",omitempty"`
}

const maxRequestBodySize = 5000

// Initialise creates the API and attaches its HTTP handler
func Initialise(server *mattrax.Server, r *mux.Router) error {
	r.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, server.Version)
	}).Methods("GET")

	r.HandleFunc("/api/settings", func(w http.ResponseWriter, r *http.Request) {
		yaml.NewEncoder(w).Encode(server.Settings.Get())
	}).Methods("GET")

	r.HandleFunc("/api/settings", func(w http.ResponseWriter, r *http.Request) {
		var cmd settings.Settings
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
		if err := yaml.NewDecoder(r.Body).Decode(&cmd); err != nil {
			res, _ := json.Marshal(Response{
				Success: false,
				Msg:     "Invalid request body",
			})
			w.Write(res)
			return
		}

		if err := cmd.Verify(); err != nil {
			res, _ := json.Marshal(Response{
				Success: false,
				Msg:     err.Error(),
			})
			w.Write(res)
			return
		}

		previousSettings := server.Settings.Get()

		server.Settings.Set(cmd)

		if previousSettings.Tenant.Name != cmd.Tenant.Name {
			if err := server.Certificates.GenerateIdentity(pkix.Name{
				CommonName: cmd.Tenant.Name + " Identity",
			}); err != nil {
				res, _ := json.Marshal(Response{
					Success: false,
					Msg:     err.Error(),
				})
				w.Write(res)
				return
			}
		}

		res, _ := json.Marshal(Response{
			Success: true,
		})
		w.Write(res)
	}).Methods("POST")

	return nil
}
