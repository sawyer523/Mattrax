package mattrax

import (
	"net/http"

	"github.com/mattrax/Mattrax/internal/api"
	"github.com/mattrax/Mattrax/internal/authentication"
	"github.com/mattrax/Mattrax/internal/mattrax/settings"
	"github.com/mattrax/Mattrax/internal/services/device"
	"github.com/mattrax/Mattrax/internal/services/policy"
	"github.com/rs/zerolog/log"

	_ "github.com/mattrax/Mattrax/internal/authentication/providers/azuread"
)

// Init starts all the Mattrax components and attaches them all to the Server
func (srv *Server) Init() {
	settingsZone, err := srv.Store.Zone("settings")
	if err != nil {
		panic(err) // TODO
	}

	ss, err := settings.NewService(settingsZone)
	if err != nil {
		panic(err) // TODO
	}
	srv.Settings = ss

	deviceZone, err := srv.Store.Zone("devices")
	if err != nil {
		panic(err) // TODO
	}

	ds, err := device.NewService(deviceZone)
	if err != nil {
		panic(err) // TODO
	}
	srv.Device = ds

	policyZone, err := srv.Store.Zone("policies")
	if err != nil {
		panic(err) // TODO
	}

	ps, err := policy.NewService(policyZone)
	if err != nil {
		panic(err) // TODO
	}
	srv.Policy = ps

	for name, prov := range authentication.Providers {
		if err := prov.Init(srv.Settings.Get().AuthProviderSettings[name]); err != nil {
			log.Fatal().Err(err).Msg("Failed initialising authentication provider!")
		}
	}

	srv.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mattrax MDM Server!"))
	})

	srv.Router.PathPrefix("/auth/{provider}").Methods(http.MethodGet, http.MethodPost).Handler(authentication.ProviderHandler(srv.Settings))
	srv.Router.PathPrefix("/auth").Name("auth").Methods(http.MethodGet).Handler(authentication.Handler())

	apiRouter := srv.Router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(api.Middleware())
	apiRouter.HandleFunc("/", api.IndexHandler(srv.Version)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/settings", api.SettingsHandler(ss)).Methods(http.MethodGet, http.MethodPatch)
	apiRouter.HandleFunc("/devices/{uuid}", api.DeviceHandler(ds)).Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)
	apiRouter.HandleFunc("/devices", api.DeviceHandler(ds)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/policies/{uuid}", api.PolicyHandler(ps)).Methods(http.MethodGet, http.MethodPatch, http.MethodDelete)
	apiRouter.HandleFunc("/policies", api.PolicyHandler(ps)).Methods(http.MethodGet)
}
