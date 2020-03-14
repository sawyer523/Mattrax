package windows

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattrax/Mattrax/internal/mattrax"
)

const maxRequestBodySize = 5000

type MDM struct {
	srv *mattrax.Server
	r   *mux.Router

	EnrollmentPolicyServiceURL    string
	EnrollmentProvisionServiceURL string
	EnrollmentManageServiceURL    string
	AuthenticationServiceURL      string
}

func Initialise(srv *mattrax.Server) (*MDM, error) {
	mdm := &MDM{
		srv: srv,
	}
	mdm.r = srv.Router.PathPrefix("/EnrollmentServer").Subrouter()
	// TODO: Middleware
	// if r.ContentLength > maxRequestBodySize {
	// 	fault := soap.NewBasicFault("s:Sender", "s:MessageFormat", "client request body too large to process")
	// 	fault.Response(w)
	// 	return
	// }
	mdm.r.HandleFunc("/Discovery.svc", mdm.DiscoveryHandler()).Methods(http.MethodGet, http.MethodPost)
	mdm.r.HandleFunc("/Policy.svc", mdm.PolicyHandler()).Name("winmdm-policy").Methods(http.MethodPost)
	mdm.r.HandleFunc("/Provision.svc", mdm.ProvisionHandler()).Name("winmdm-provision").Methods(http.MethodPost)
	mdm.r.HandleFunc("/Manage.svc", mdm.ManageHandler()).Name("winmdm-manage").Methods(http.MethodPost)

	// TODO: Replace with UI
	mdm.r.HandleFunc("/ToS", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>MDM Consent</title></head><body><h3>MDM Consent</h3><button onClick="acceptBtn()">Accept</button><button onClick="denyBtn()">Reject</button><script>
function acceptBtn() {
	var urlParams = new URLSearchParams(window.location.search);

	if (!urlParams.has('redirect_uri')) {
		alert('Redirect url not found. Did you open this in your broswer?');
	} else {
		window.location = urlParams.get('redirect_uri') + "?IsAccepted=true&OpaqueBlob=TODOCustomDataFromAzureAD";
	}
}
function denyBtn() {
	var urlParams = new URLSearchParams(window.location.search);

	if (!urlParams.has('redirect_uri')) {
		alert('Redirect url not found. Did you open this in your broswer?');
	} else {
		window.location = urlParams.get('redirect_uri') + "?IsAccepted=false&error=access_denied&error_description=Access%20is%20denied%2E";
	}
}</script></body></html>`)
	}).Methods(http.MethodGet)

	url, err := mdm.r.Get("winmdm-policy").URL()
	if err != nil {
		panic(err) // TODO
	}
	url.Scheme = "https"
	url.Host = srv.Config.Domain
	mdm.EnrollmentPolicyServiceURL = url.String()

	url, err = mdm.r.Get("winmdm-provision").URL()
	if err != nil {
		panic(err) // TODO
	}
	url.Scheme = "https"
	url.Host = srv.Config.Domain
	mdm.EnrollmentProvisionServiceURL = url.String()

	url, err = mdm.r.Get("winmdm-manage").URL()
	if err != nil {
		panic(err) // TODO
	}
	url.Scheme = "https"
	url.Host = srv.Config.Domain
	mdm.EnrollmentManageServiceURL = url.String()

	url, err = srv.Router.Get("auth").URL()
	if err != nil {
		panic(err) // TODO
	}
	url.Scheme = "https"
	url.Host = srv.Config.Domain
	mdm.AuthenticationServiceURL = url.String()

	return mdm, nil
}
