package portals

import (
	"fmt"
	"net/http"
)

func AzureTOSHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization") // AzureAD JWT
		// fmt.Println(authHeader)
		if authHeader == "" {
			fmt.Fprintf(w, `<html>
		<head>
			<title>MDM Concent</title>
		</head>
		<body>
			<h3>Failed Authorization</h3>
		</body>
		</html>`)
			return
		}

		// TODO
		fmt.Fprintf(w, `<html>
		<head>
			<title>MDM Concent</title>
		</head>
		<body>
			<h3>MDM Concent</h3>
			<button onClick="acceptBtn()">Accept</button>
			<button onClick="denyBtn()">Reject</button>
			<script>
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
			}
			</script>
		</body>
		</html>`)
	}
}
