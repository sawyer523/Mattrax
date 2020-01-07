package portals

import (
	"fmt"
	"net/http"
)

func FederatedLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
