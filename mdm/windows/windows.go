package windows

import (
	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/mdm/windows/protocol"
)

// Init initialises the Windows MDM components
func Init(server mattrax.Server, r *mux.Router) error {
	r.Path("/EnrollmentServer/Discovery.svc").Methods("GET").HandlerFunc(protocol.Discover(server))
	r.Path("/EnrollmentServer/Discovery.svc").Methods("POST").HandlerFunc(protocol.Discovery(server))
	r.Path("/EnrollmentServer/Policy.svc").Methods("POST").HandlerFunc(protocol.Policy(server))

	return nil
}
