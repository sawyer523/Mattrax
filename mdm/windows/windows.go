package windows

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	enrolldiscovery "github.com/mattrax/Mattrax/mdm/windows/protocol/enroll_discovery"
	enrollpolicy "github.com/mattrax/Mattrax/mdm/windows/protocol/enroll_policy"
	enrollprovision "github.com/mattrax/Mattrax/mdm/windows/protocol/enroll_provision"
)

// Init initialises the Windows MDM components
func Init(server mattrax.Server, r *mux.Router) error {
	r.Path("/EnrollmentServer/Discovery.svc").Methods("GET").HandlerFunc(verifyUserAgent("ENROLLClient", enrolldiscovery.GETHandler(server)))
	r.Path("/EnrollmentServer/Discovery.svc").Methods("POST").HandlerFunc(verifyUserAgent("ENROLLClient", verifySoapRequest(enrolldiscovery.Handler(server))))
	r.Path("/EnrollmentServer/Policy.svc").Methods("POST").HandlerFunc(verifyUserAgent("ENROLLClient", verifySoapRequest(enrollpolicy.Handler(server))))
	r.Path("/EnrollmentServer/Enrollment.svc").Methods("POST").HandlerFunc(verifyUserAgent("ENROLLClient", verifySoapRequest(enrollprovision.Handler(server))))

	// TODO: Cleanup below + move to own package + Configurable Internal ones + Allow Custom External URL
	r.Path("/EnrollmentServer/Authenticate").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html>
		<head>
			<title>MDM Federated Login</title>
		</head>
		<body>
			<h3>MDM Federated Login</h3>
			<form method="post" action="ms-app://windows.immersivecontrolpanel">
				<p><input type="hidden" name="wresult" value="TODOSpecialTokenWhichVerifiesAuth" /></p>
				<input type="submit" value="Login" />
			</form>
		</body>
		</html>`)
	})

	r.Path("/EnrollmentServer/ToS").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html>
		<head>
			<title>MDM Concent</title>
		</head>
		<body>
			<h3>MDM Concent</h3>
			<button onClick="acceptBtn()">Accept</button>
			<script>
			function acceptBtn() {
				var urlParams = new URLSearchParams(window.location.search);

				if (!urlParams.has('redirect_uri')) {
					alert('Redirect url not found. Did you open this in your broswer?');
				} else {
					window.location = urlParams.get('redirect_uri') + "?IsAccepted=true&OpaqueBlob=TODOCustomDataFromAzureAD";
				}
			}
			</script>
		</body>
		</html>`)
	})

	return nil
}
