package windows

import (
	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	enrolldiscovery "github.com/mattrax/Mattrax/mdm/windows/protocol/enroll_discovery"
	enrollpolicy "github.com/mattrax/Mattrax/mdm/windows/protocol/enroll_policy"
	enrollprovision "github.com/mattrax/Mattrax/mdm/windows/protocol/enroll_provision"
	"github.com/mattrax/Mattrax/mdm/windows/protocol/portals"
)

// Init initialises the Windows MDM components
func Init(server *mattrax.Server, r *mux.Router) error {
	r.Path("/EnrollmentServer/Discovery.svc").Methods("GET").HandlerFunc(defaultHeaders("ENROLLClient", enrolldiscovery.GETHandler()))
	r.Path("/EnrollmentServer/Discovery.svc").Methods("POST").HandlerFunc(defaultHeaders("ENROLLClient", enrolldiscovery.Handler(server)))
	r.Path("/EnrollmentServer/Policy.svc").Methods("POST").HandlerFunc(defaultHeaders("ENROLLClient", enrollpolicy.Handler(server)))
	r.Path("/EnrollmentServer/Enrollment.svc").Methods("POST").HandlerFunc(defaultHeaders("ENROLLClient", enrollprovision.Handler(server)))
	r.Path("/EnrollmentServer/Authenticate").Methods("GET").HandlerFunc(portals.FederatedLoginHandler())
	r.Path("/EnrollmentServer/ToS").Methods("GET").HandlerFunc(portals.AzureTOSHandler())

	return nil
}
