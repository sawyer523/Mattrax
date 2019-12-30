package mdm

import (
	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/mdm/windows"
)

// Initialise starts the MDM services and attaches thier HTTP handlers to the router
func Initialise(server *mattrax.Server, r *mux.Router) error {
	if err := windows.Init(server, r); err != nil {
		return err
	}

	return nil
}

// Deinitialise runs before the server shutsdown and it should cleanup and stop running services
func Deinitialise() error {
	// TODO
	return nil
}
