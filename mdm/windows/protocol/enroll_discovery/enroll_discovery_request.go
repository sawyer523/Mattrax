package enrolldiscovery

import (
	"strings"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	generic "github.com/mattrax/Mattrax/mdm/windows/protocol/generic"
	wsettings "github.com/mattrax/Mattrax/mdm/windows/settings"
	"github.com/mattrax/Mattrax/pkg/xml"
	"github.com/pkg/errors"
)

type AuthPolicies struct {
	AuthPolicy []string `xml:"AuthPolicy"`
}

type RequestBody struct {
	EmailAddress       string       `xml:"EmailAddress"`
	RequestVersion     string       `xml:"RequestVersion"`
	DeviceType         string       `xml:"DeviceType"`
	ApplicationVersion string       `xml:"ApplicationVersion"`
	OSEdition          string       `xml:"OSEdition"`
	AuthPolicies       AuthPolicies `xml:"AuthPolicies"`
}

type Request struct {
	XMLName xml.Name       `xml:"s:Envelope"`
	Header  generic.Header `xml:"s:Header"`
	Body    RequestBody    `xml:"s:Body>Discover>request"`
}

func (cmd Request) Verify(config mattrax.Config, settings types.Settings) error {
	/* Verify Structure: Lightwieght structure and datatype checks */
	if err := cmd.Header.VerifyStructure("http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover", false); err != nil {
		return err
	}

	if cmd.Body.EmailAddress == "" {
		return errors.New("empty EmailAddress")
	} else if !types.ValidEmail.MatchString(cmd.Body.EmailAddress) {
		return errors.New("invalid EmailAddress '" + cmd.Body.EmailAddress + "'")
	}

	if cmd.Body.RequestVersion != "4.0" {
		return errors.New("invalid RequestVersion '" + cmd.Body.RequestVersion + "' only '4.0' is supported")
	}

	if cmd.Body.DeviceType != "" && cmd.Body.DeviceType != "CIMClient_Windows" {
		return errors.New("invalid DeviceType '" + cmd.Body.DeviceType + "' only 'CIMClient_Windows' are supported")
	}

	if len(cmd.Body.AuthPolicies.AuthPolicy) > 0 {
		for _, authPolicy := range cmd.Body.AuthPolicies.AuthPolicy {
			if !(authPolicy == "Federated" || authPolicy == "OnPremise" || authPolicy == "Certificate") {
				return errors.New("invalid supported AuthPolicy '" + authPolicy + "' only 'Federated' or 'OnPremise' or 'Certificate' is supported")
			}
		}
	}

	/* Verify Context: Expensive checks against the server's DB */
	if err := cmd.Header.VerifyContext(config); err != nil {
		return err
	}

	// Verify valid email domain // TODO: Redo this logic to be cleaner
	email := strings.Split(cmd.Body.EmailAddress, "@")
	if len(email) != 2 {
		// This should NEVER be called because the email is verified in VerifyStructure
		return errors.New("invalid EmailAddress '" + cmd.Body.EmailAddress + "'")
	}

	for _, domain := range settings.ManagedDomains {
		if email[1] == domain {
			return nil
		}
	}

	return errors.New("the request failed verification") // TODO: Better warning about ManagedDomains
}

func (cmd Request) IsAuthPolicySupport(authPolicy wsettings.AuthPolicy) bool {
	for _, ap := range cmd.Body.AuthPolicies.AuthPolicy {
		if ap == "OnPremise" && authPolicy == wsettings.AuthPolicyOnPremise {
			return true
		}
		if ap == "Federated" && authPolicy == wsettings.AuthPolicyFederated {
			return true
		}
		if ap == "Certificate" && authPolicy == wsettings.AuthPolicyCertificate {
			return true
		}
	}

	return false
}
