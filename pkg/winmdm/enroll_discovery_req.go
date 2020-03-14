package winmdm

import (
	"net/url"

	"github.com/mattrax/Mattrax/pkg/types"
	"github.com/mattrax/Mattrax/pkg/winmdm/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
)

// DiscoveryRequest is the SOAP Envelope sent from the client during the Discovery phase
type DiscoveryRequest struct {
	XMLName xml.Name             `xml:"s:Envelope"`
	Header  soap.Header          `xml:"s:Header"`
	Body    DiscoveryRequestBody `xml:"s:Body>Discover>request"`
}

// Verify checks the DiscoveryRequest contains valid values and if not it returns a soap fault
func (cmd DiscoveryRequest) Verify(desiredAuthPolicy string) (bool, soap.FaultEnvelop) {
	if cmd.Header.Action != "http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover" {
		return false, soap.NewBasicFault("s:Sender", "a:ActionMismatch", "Client request body is invalid")
	}

	if url, err := url.ParseRequestURI(cmd.Header.To); err != nil || url.Scheme != "https" {
		return false, soap.NewEnrollmentFault("s:Sender", "a:EndpointUnavailable", "Invalid To address", "EnrollmentServer", "EnrollmentInternalServiceError", "")
	}

	if cmd.Body.EmailAddress == "" || !types.ValidEmail.MatchString(cmd.Body.EmailAddress) {
		return false, soap.NewBasicFault("s:Sender", "s:MessageFormat", "Client request body missing or invalid email address")
	}

	// Note: Intune doesn't verify this but it is here as it should prevent issues after enrollment if an unsupported device was enrolled
	if cmd.Body.RequestVersion != "4.0" {
		return false, soap.NewEnrollmentFault("s:Sender", "a:InternalServiceFault", "The discovery request version "+cmd.Body.RequestVersion+" is not supported", "DeviceNotSupported", "unsupported discovery request version", "")
	}

	if cmd.Body.DeviceType != "CIMClient_Windows" {
		return false, soap.NewEnrollmentFault("s:Sender", "a:InternalServiceFault", "Device type "+cmd.Body.DeviceType+" is not supported", "DeviceNotSupported", "unsupported device type", "")
	}

	// Note: Intune disregards what the device supports and returns the AuthPolicy it desires but this was not done here
	if !cmd.Body.AuthPolicies.IsAuthPolicySupported(desiredAuthPolicy) {
		return false, soap.NewEnrollmentFault("s:Sender", "a:InternalServiceFault", desiredAuthPolicy+" auth policy is not supported by your device by required for enrollment", "DeviceNotSupported", "unsupported auth policy", "")
	}

	return true, soap.FaultEnvelop{}
}

// DiscoveryRequestBody contains body of the request
type DiscoveryRequestBody struct {
	EmailAddress       string       `xml:"EmailAddress"`
	RequestVersion     string       `xml:"RequestVersion"`
	DeviceType         string       `xml:"DeviceType"`
	ApplicationVersion string       `xml:"ApplicationVersion"`
	OSEdition          string       `xml:"OSEdition"`
	AuthPolicies       AuthPolicies `xml:"AuthPolicies"`
}

// AuthPolicies is an array of the supported AuthPolicies
type AuthPolicies struct {
	AuthPolicies []string `xml:"AuthPolicy"`
}

// IsAuthPolicySupported checks whether an AuthPolicy exists in the AuthPolicies array
func (authPolicies AuthPolicies) IsAuthPolicySupported(authPolicyStr string) bool {
	for _, ap := range authPolicies.AuthPolicies {
		if ap == authPolicyStr {
			return true
		}
	}
	return false
}
