package enrolldiscovery

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
)

// Add Tests For:
// - Too big request body
// - Test multiple configurations (AuthPolicies, Values set and not set)
// - Returns fault for
//		- Invalid Action
//		- Invalid To adrress + Verify allowed address's include config.Domain and settings.ManagedDomains

var server = &mattrax.Server{
	// Settings: types.Settings{
	// 	Windows: wsettings.Settings{
	// 		DeploymentType:      wsettings.DeploymentStandalone,
	// 		FederationPortalURL: "",
	// 	},
	// },
} // TOOD: Mock fill in -> Using custom DB path

func TestDiscoveryGETResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "/EnrollServer/Discovery.svc", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	GETHandler(server)(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("DiscoveryGET: returned an incorrect status code: got %v want %v", status, http.StatusOK)
	}

	if res.Body.Len() != 0 {
		t.Errorf("DiscoveryGET: didn't return a blank body: got '%v' want '%v'", res.Body.String(), "")
	}
}

func TestDiscoveryPOSTMissingRequestBody(t *testing.T) {
	// TODO: RequestBody

	req, err := http.NewRequest("POST", "/EnrollServer/Discovery.svc", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	Handler(server)(res, req)

	if status := res.Code; status != http.StatusInternalServerError {
		t.Errorf("TestDiscoveryPOSTMissingRequestBody: returned an incorrect status code: got %v want %v", status, http.StatusInternalServerError)
	}

	if res.Body.Len() == 0 {
		t.Errorf("TestDiscoveryPOSTMissingRequestBody: missing response body'")
		return
	}

	var fault soap.FaultEnvelop
	if err := xml.NewDecoder(res.Body).Decode(&fault); err != nil {
		t.Errorf("TestDiscoveryPOSTMissingRequestBody: failed to parse fault: got error %v", err)
	}

	if fault.Value != "s:Sender" {
		t.Errorf("TestDiscoveryPOSTMissingRequestBody: fault cause was incorrect: got %v want %v", fault.Value, "s:Sender")
	}

	// Although the MDM spec allows this to be multiple possible value this is done to ensure a logical and supported value is used.
	if fault.Subcode != "s:MessageFormat" {
		t.Errorf("TestDiscoveryPOSTMissingRequestBody: fault subcode value was incorrect: got %v want %v", fault.Subcode, "s:MessageFormat")
	}

	if fault.Reason.Text == "" {
		t.Errorf("TestDiscoveryPOSTMissingRequestBody: fault reason text was not given: got %v", fault.Reason.Text)
	}
}

func TestDiscoveryPOST(t *testing.T) {
	req, err := http.NewRequest("POST", "/EnrollServer/Discovery.svc", nil)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()
	Handler(server)(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("TestDiscoveryPOST: returned an incorrect status code: got %v want %v", status, http.StatusOK)
	}

	if res.Body.Len() == 0 {
		t.Errorf("TestDiscoveryPOST: missing response body'")
		return
	}

	var cmd Response
	var cmdRes = cmd.Body.DiscoverResponse
	if err := xml.NewDecoder(res.Body).Decode(&cmd); err != nil {
		t.Errorf("TestDiscoveryPOST: failed to parse cmd: got error %v", err)
	}

	if cmdRes.AuthPolicy != "OnPremise" && cmdRes.AuthPolicy != "Federated" && cmdRes.AuthPolicy != "Certificate" {
		t.Errorf("TestDiscoveryPOST: returned an invalid auth policy: got %v want OnPremise, Federated or Certificate", cmdRes.AuthPolicy)
		return
	}

	if cmdRes.EnrollmentVersion != "4.0" {
		t.Errorf("TestDiscoveryPOST: returned unsupported enrollment version: got %v want %v", cmdRes.EnrollmentVersion, "4.0")
		return
	}

	if cmdRes.EnrollmentPolicyServiceURL != "" {
		t.Errorf("TestDiscoveryPOST: failed to return an enrollment policy service url")
		return
	} else if url, err := url.ParseRequestURI(cmdRes.EnrollmentPolicyServiceURL); err != nil {
		t.Errorf("TestDiscoveryPOST: invalid enrollment policy service url: url '%v' error %v", cmdRes.EnrollmentServiceURL, err)
		return
	} else if url.Scheme != "https" {
		t.Errorf("TestDiscoveryPOST: invalid enrollment policy service url scheme: got %v want %v", url.Scheme, "https")
	}

	if cmdRes.EnrollmentServiceURL != "" {
		t.Errorf("TestDiscoveryPOST: failed to return an enrollment service url")
		return
	} else if url, err := url.ParseRequestURI(cmdRes.EnrollmentServiceURL); err != nil {
		t.Errorf("TestDiscoveryPOST: invalid enrollment service url: url '%v' error %v", cmdRes.EnrollmentServiceURL, err)
		return
	} else if url.Scheme != "https" {
		t.Errorf("TestDiscoveryPOST: invalid enrollment service url scheme: got %v want %v", url.Scheme, "https")
	}

	if cmdRes.AuthPolicy == "Federated" {
		if cmdRes.AuthenticationServiceURL != "" {
			t.Errorf("TestDiscoveryPOST: failed to return an authentication service url")
			return
		} else if url, err := url.ParseRequestURI(cmdRes.AuthenticationServiceURL); err != nil {
			t.Errorf("TestDiscoveryPOST: invalid authentication service url: url '%v' error %v", cmdRes.AuthenticationServiceURL, err)
			return
		} else if url.Scheme != "https" {
			t.Errorf("TestDiscoveryPOST: invalid authentication service url scheme: got %v want %v", url.Scheme, "https")
		}
	}
}
