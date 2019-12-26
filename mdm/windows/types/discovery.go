package wtypes

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/mattrax/Mattrax/pkg/xml"
)

// MdeDiscoveryRequest contains the payload sent by the client to discover the enrollment services
type MdeDiscoveryRequest struct {
	XMLName xml.Name `xml:"s:Envelope"`
	Header  struct {
		Action    string `xml:"a:Action"`
		MessageID string `xml:"a:MessageID"`
		ReplyTo   struct {
			Address string `xml:"a:Address"`
		} `xml:"a:ReplyTo"`
		To string `xml:"a:To"`
	} `xml:"s:Header"`
	Body struct {
		Discover struct {
			Request struct {
				EmailAddress       string `xml:"EmailAddress"`
				RequestVersion     string `xml:"RequestVersion"`
				DeviceType         string `xml:"DeviceType"`
				ApplicationVersion string `xml:"ApplicationVersion"`
				OSEdition          string `xml:"OSEdition"`
				AuthPolicies       struct {
					AuthPolicy []string `xml:"AuthPolicy"`
				} `xml:"AuthPolicies"`
			} `xml:"request"`
		} `xml:"Discover"`
	} `xml:"s:Body"`
}

// validEmail is a regex used to verify an email is valid
var validEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// validMessageID is a regex used to verify an MessageID is valid
var validMessageID = regexp.MustCompile(`^[a-zA-Z0-9:\-]+$`)

// VerifyStructure checks that the cmd contains a valid structure.
// This is done before more resource intensive checks to reduce the DoS vector.
func (cmd MdeDiscoveryRequest) VerifyStructure() error {
	if cmd.Header.Action == "" {
		return errors.New("invalid MdeDiscoveryRequest: empty Action")
	} else if cmd.Header.Action != "http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover" {
		return errors.New("invalid MdeDiscoveryRequest: invalid Action '" + cmd.Header.Action + "'")
	}

	if cmd.Header.MessageID == "" {
		return errors.New("invalid MdeDiscoveryRequest: empty MessageID")
	} else if !strings.HasPrefix(cmd.Header.MessageID, "urn:uuid:") {
		return errors.New("invalid MdeDiscoveryRequest: invalid MessageID prefix '" + cmd.Header.MessageID + "'")
	} else if !validMessageID.MatchString(cmd.Header.MessageID) {
		return errors.New("invalid MdeDiscoveryRequest: invalid characters in MessageID '" + cmd.Header.MessageID + "'")
	}

	if cmd.Header.To == "" {
		return errors.New("invalid MdeDiscoveryRequest: empty To")
	} else if _, err := url.ParseRequestURI(cmd.Header.To); err != nil {
		return errors.New("invalid MdeDiscoveryRequest: invalid To '" + cmd.Header.To + "'")
	}

	if cmd.Body.Discover.Request.EmailAddress == "" {
		return errors.New("invalid MdeDiscoveryRequest: empty EmailAddress")
	} else if !validEmail.MatchString(cmd.Body.Discover.Request.EmailAddress) {
		return errors.New("invalid MdeDiscoveryRequest: invalid EmailAddress '" + cmd.Body.Discover.Request.EmailAddress + "'")
	}

	if cmd.Body.Discover.Request.RequestVersion != "4.0" {
		return errors.New("invalid MdeDiscoveryRequest: invalid RequestVersion '" + cmd.Body.Discover.Request.RequestVersion + "' only '4.0' is supported")
	}

	if cmd.Body.Discover.Request.DeviceType != "" && cmd.Body.Discover.Request.DeviceType != "CIMClient_Windows" {
		return errors.New("invalid MdeDiscoveryRequest: invalid DeviceType '" + cmd.Body.Discover.Request.DeviceType + "' only 'CIMClient_Windows' are supported")
	}

	if len(cmd.Body.Discover.Request.AuthPolicies.AuthPolicy) > 0 {
		for _, authPolicy := range cmd.Body.Discover.Request.AuthPolicies.AuthPolicy {
			if !(authPolicy == "Federated" || authPolicy == "OnPremise" || authPolicy == "Certificate") {
				return errors.New("invalid MdeDiscoveryRequest: invalid supported AuthPolicy '" + authPolicy + "' only 'Federated' or 'OnPremise' or 'Certificate' is supported")
			}
		}
	}

	return nil
}

// VerifyContext checks that the cmd against the expected values.
// This operation is more expensive so should always be done after VerifyStructure.
func (cmd MdeDiscoveryRequest) VerifyContext(config mattrax.Config, settings types.Settings) error {
	// Verify valid To address
	if toAddrRaw, err := url.Parse(cmd.Header.To); err != nil {
		// This should NEVER be called because the url is verified in VerifyStructure
		return errors.New("invalid MdeDiscoveryRequest: invalid To '" + cmd.Header.To + "'")
	} else if !(strings.ToLower(toAddrRaw.Hostname()) == strings.ToLower(config.PrimaryDomain) || strings.ToLower(toAddrRaw.Hostname()) == strings.ToLower(config.WindowsDiscoveryDomain)) {
		return errors.New("invalid MdeDiscoveryRequest: this requested server ('" + strings.ToLower(toAddrRaw.Hostname()) + "') isn't this server ('" + config.WindowsDiscoveryDomain + "' or '" + config.PrimaryDomain + "')")
	}

	// Verify valid email domain
	email := strings.Split(cmd.Body.Discover.Request.EmailAddress, "@")
	if len(email) != 2 {
		// This should NEVER be called because the email is verified in VerifyStructure
		return errors.New("invalid MdeDiscoveryRequest: invalid EmailAddress '" + cmd.Body.Discover.Request.EmailAddress + "'")
	}

	for _, domain := range settings.ManagedDomains {
		if email[1] == domain {
			return nil
		}
	}

	return errors.New("invalid MdeDiscoveryRequest: the request failed verification") // TODO: Better warning about ManagedDomains
}

func (cmd MdeDiscoveryRequest) IsAuthPolicySupport(authPolicy types.AuthPolicy) bool {
	for _, ap := range cmd.Body.Discover.Request.AuthPolicies.AuthPolicy {
		if ap == "OnPremise" && authPolicy == types.AuthPolicyOnPremise {
			return true
		}
		if ap == "Federated" && authPolicy == types.AuthPolicyFederated {
			return true
		}
		if ap == "Certificate" && authPolicy == types.AuthPolicyCertificate {
			return true
		}
	}

	return false
}

// DiscoverResponse is the payload within the body of the response
type DiscoverResponse struct {
	AuthPolicy                 string `xml:"DiscoverResult>AuthPolicy"`
	EnrollmentVersion          string `xml:"DiscoverResult>EnrollmentVersion"`
	EnrollmentPolicyServiceURL string `xml:"DiscoverResult>EnrollmentPolicyServiceUrl"`
	EnrollmentServiceURL       string `xml:"DiscoverResult>EnrollmentServiceUrl"`
	AuthenticationServiceUrl   string `xml:"DiscoverResult>AuthenticationServiceUrl,omitempty"`
}

// MdeDiscoveryResponseBody is the body tag of the response
type MdeDiscoveryResponseBody struct {
	NamespaceXSI     string           `xml:"xmlns:xsi,attr"`
	NamespaceXSD     string           `xml:"xmlns:xsd,attr"`
	DiscoverResponse DiscoverResponse `xml:"http://schemas.microsoft.com/windows/management/2012/01/enrollment DiscoverResponse"`
}

// MdeDiscoveryResponseEnvelope contains the payload sent from the server to the client telling the client of its capabilities and endpoints for MDE
type MdeDiscoveryResponseEnvelope struct {
	XMLName          xml.Name                 `xml:"s:Envelope"`
	NamespaceS       string                   `xml:"xmlns:s,attr"`
	NamespaceA       string                   `xml:"xmlns:a,attr"`
	HeaderAction     MustUnderstand           `xml:"s:Header>a:Action"`
	HeaderActivityID string                   `xml:"s:Header>a:ActivityID"`
	HeaderRelatesTo  string                   `xml:"s:Header>a:RelatesTo"`
	Body             MdeDiscoveryResponseBody `xml:"s:Body"`
}
