package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/api"
	"github.com/mattrax/Mattrax/internal/boltdb"
	"github.com/mattrax/Mattrax/internal/middleware"
	"github.com/mattrax/Mattrax/mdm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Allow non zero return code while retaining defer statement execution
	// Also handle panic statements which stopped working
	returnCode := 0
	defer func() {
		if errRaw := recover(); errRaw != nil {
			if err, ok := errRaw.(error); ok {
				log.Error().Err(err).Msg("Panic!")
			} else {
				log.Error().Interface("error", errRaw).Msg("Panic!")
			}
			os.Exit(1)
		} else {
			os.Exit(returnCode)
		}
	}()

	// Parse and verify command line flags
	config := mattrax.Config{}
	p := arg.MustParse(&config)
	config.Verify(p)

	// Create server
	server := &mattrax.Server{
		Version: mattrax.Version,
		Config:  config,
	}

	// Initialise logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.DevelopmentMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Initialise datastore
	if err := boltdb.Initialise(server); err != nil {
		log.Error().Str("dbpath", config.DBPath).Err(err).Msg("Error initialising the datastore!")
		returnCode = 1
		return
	}
	defer func() {
		if err := boltdb.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing the datastore!")
			returnCode = 1
		}
	}()

	// Initialise router and HTTP server
	r := mux.NewRouter()
	httpSrv := &http.Server{
		Addr:         ":" + strconv.Itoa(config.Port),
		Handler:      middleware.Global(r, config.DevelopmentMode),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
			MinVersion: tls.VersionTLS12,
		},
		// FUTURE: ErrorLog: {MAKE COMPATIBLE},
	}

	// Initialise MDM protocols
	if err := mdm.Initialise(server, r); err != nil {
		log.Error().Err(err).Msg("Error initialising the MDM protocols!")
		returnCode = 1
		return
	}
	defer func() {
		if err := mdm.Deinitialise(); err != nil {
			log.Error().Err(err).Msg("Error deinitialising the MDM protocols!")
			returnCode = 1
		}
	}()

	// Initialise API
	if err := api.Initialise(server, r); err != nil {
		log.Error().Err(err).Msg("Error initialising the API!")
		returnCode = 1
		return
	}

	// Create gracefull shutdown channel
	done := make(chan os.Signal, 1)

	// Start the HTTP server
	go func() {
		if err := httpSrv.ListenAndServeTLS(config.CertFile, config.KeyFile); err != nil && err != http.ErrServerClosed {
			log.Error().Int("port", config.Port).Str("certfile", config.CertFile).Str("keyfile", config.KeyFile).Err(err).Msg("Error with the webserver!")
			done <- syscall.Signal(-1)
		}
	}()

	// Everything Good
	if config.DevelopmentMode {
		log.Warn().Msg("Development Mode Enabled. DO NOT use the server in production!")
	}
	log.Info().Str("domain", config.Domain).Int("port", config.Port).Msg("Initialised Mattrax MDM Server!")

	// Upon shutdown request gracefully close the HTTP server
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-done

	timeout := 15 * time.Second
	if config.DevelopmentMode {
		timeout = 2 * time.Second
	}

	if sig.String() != "signal -1" {
		log.Info().Msg("Please wait while the server is gracefully shutdown. It will timeout after " + timeout.String() + " seconds if unable to complete.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := httpSrv.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error shutting down the webserver!")
	}
}
