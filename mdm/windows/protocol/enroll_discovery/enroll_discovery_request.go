package enrolldiscovery

import (
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
)

// Request contains the SOAP request Envelope
type Request struct {
	XMLName xml.Name    `xml:"s:Envelope"`
	Header  soap.Header `xml:"s:Header"`
	Body    RequestBody `xml:"s:Body>Discover>request"`
}

// RequestBody contains the body of the SOAP Envelope
type RequestBody struct {
	EmailAddress       string       `xml:"EmailAddress"`
	RequestVersion     string       `xml:"RequestVersion"`
	DeviceType         string       `xml:"DeviceType"`
	ApplicationVersion string       `xml:"ApplicationVersion"`
	OSEdition          string       `xml:"OSEdition"`
	AuthPolicies       AuthPolicies `xml:"AuthPolicies"`
}

// AuthPolicies contains the array of supported AuthPolicies
type AuthPolicies struct {
	AuthPolicies []string `xml:"AuthPolicy"`
}

// IsAuthPolicySupported checks the AuthPolicies array for the existant of an AuthPolicy
func (authPolicies AuthPolicies) IsAuthPolicySupported(authPolicyStr string) bool {
	for _, ap := range authPolicies.AuthPolicies {
		if ap == authPolicyStr {
			return true
		}
	}

	return false
}
