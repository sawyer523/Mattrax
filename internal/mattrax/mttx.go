package mattrax

import (
	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
	"github.com/mattrax/Mattrax/internal/datastore"
	datastoreinit "github.com/mattrax/Mattrax/internal/datastore/init"
	"github.com/mattrax/Mattrax/internal/mattrax/settings"
	"github.com/mattrax/Mattrax/internal/services/device"
	"github.com/mattrax/Mattrax/internal/services/policy"
	"github.com/rs/zerolog/log"
)

// Version contains the Mattrax server version
// This varible's correct value is injected at build time
var Version string = "0.0.0-development"

// Server holds the global server state
type Server struct {
	Version  string
	Config   Config
	Store    datastore.Store // TODO: Make unexported if possible
	Router   *mux.Router     // TODO: Make unexported if possible
	Settings *settings.Service
	Device   *device.Service
	Policy   *policy.Service
}

// Close gracefully stops and cleans up all the Mattrax server components
func (srv *Server) Close() {
	if err := srv.Store.Close(); err != nil {
		log.Fatal().Err(err).Msg("Error closing datastore connection!")
	}
}

// NewServer initialises and returns a new Mattrax server
func NewServer() *Server {
	srv := &Server{
		Version: Version,
		Config:  Config{},
	}
	srv.Config.Verify(arg.MustParse(&srv.Config))
	store, err := datastoreinit.Init(srv.Config.Database)
	if err != nil {
		log.Fatal().Str("connURL", srv.Config.Database).Err(err).Msg("Datastore initialisation error!")
	}
	srv.Store = store
	srv.Router = mux.NewRouter()
	srv.Init()
	return srv
}
