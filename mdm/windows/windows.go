package windows

import (
	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/mdm/windows/azuread"
	enrolldiscovery "github.com/mattrax/Mattrax/mdm/windows/protocol/enroll_discovery"
	enrollpolicy "github.com/mattrax/Mattrax/mdm/windows/protocol/enroll_policy"
	enrollprovision "github.com/mattrax/Mattrax/mdm/windows/protocol/enroll_provision"
	"github.com/mattrax/Mattrax/mdm/windows/protocol/portals"
)

// MDM is the global state container for Windows MDM
type MDM struct {
	AAD azuread.Service
}

// Init initialises the Windows MDM components
func Init(server *mattrax.Server, r *mux.Router) (MDM, error) {
	mdm := MDM{
		AAD: azuread.Service{},
	}

	// fmt.Println(mdm.AAD)
	// mdm.AAD.RefreshAccessToken(server)
	// fmt.Println(mdm.AAD)

	// TODO: expose mdm to handlers
	r.Path("/EnrollmentServer/Discovery.svc").Methods("GET").HandlerFunc(defaultHeaders(enrolldiscovery.GETHandler(server)))
	r.Path("/EnrollmentServer/Discovery.svc").Methods("POST").HandlerFunc(defaultHeaders(enrolldiscovery.Handler(server)))
	r.Path("/EnrollmentServer/Policy.svc").Methods("POST").HandlerFunc(defaultHeaders(enrollpolicy.Handler(server)))
	r.Path("/EnrollmentServer/Enrollment.svc").Methods("POST").HandlerFunc(defaultHeaders(enrollprovision.Handler(server)))
	r.Path("/ManagementServer/Manage.svc").Methods("POST").HandlerFunc(defaultHeaders(enrollprovision.Handler(server)))
	r.Path("/EnrollmentServer/Authenticate").Methods("GET").HandlerFunc(portals.FederatedLoginHandler())
	r.Path("/EnrollmentServer/ToS").Methods("GET").HandlerFunc(portals.AzureTOSHandler())

	return mdm, nil
}
