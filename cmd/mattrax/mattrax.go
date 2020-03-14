package main

import (
	"os"

	"github.com/mattrax/Mattrax/internal/http"
	"github.com/mattrax/Mattrax/internal/mattrax"
	"github.com/mattrax/Mattrax/mdm/windows"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	/* Zerologger Init: This will likely be deprecated in the future and wrapped by elog */
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if true { // TODO: When in development mode only
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	/* END Zerologger init */

	srv := mattrax.NewServer()
	defer srv.Close()

	// TODO: Move this into the Mattrax package
	_, err := windows.Initialise(srv)
	if err != nil {
		panic(err) // TODO
	}

	http.Serve(srv)
}
