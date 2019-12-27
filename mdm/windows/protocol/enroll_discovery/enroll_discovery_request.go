package enrolldiscovery

import (
	"strings"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	generic "github.com/mattrax/Mattrax/mdm/windows/protocol/generic"
	"github.com/mattrax/Mattrax/pkg/xml"
	"github.com/pkg/errors"
)

type Request struct {
	XMLName xml.Name       `xml:"s:Envelope"`
	Header  generic.Header `xml:"s:Header"`
	Body    struct {
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

func (cmd Request) Verify(config mattrax.Config, settings types.Settings) error {
	/* Verify Structure: Lightwieght structure and datatype checks */
	if err := cmd.Header.VerifyStructure("http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover", false); err != nil {
		return err
	}

	if cmd.Body.Discover.Request.EmailAddress == "" {
		return errors.New("empty EmailAddress")
	} else if !types.ValidEmail.MatchString(cmd.Body.Discover.Request.EmailAddress) {
		return errors.New("invalid EmailAddress '" + cmd.Body.Discover.Request.EmailAddress + "'")
	}

	if cmd.Body.Discover.Request.RequestVersion != "4.0" {
		return errors.New("invalid RequestVersion '" + cmd.Body.Discover.Request.RequestVersion + "' only '4.0' is supported")
	}

	if cmd.Body.Discover.Request.DeviceType != "" && cmd.Body.Discover.Request.DeviceType != "CIMClient_Windows" {
		return errors.New("invalid DeviceType '" + cmd.Body.Discover.Request.DeviceType + "' only 'CIMClient_Windows' are supported")
	}

	if len(cmd.Body.Discover.Request.AuthPolicies.AuthPolicy) > 0 {
		for _, authPolicy := range cmd.Body.Discover.Request.AuthPolicies.AuthPolicy {
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
	email := strings.Split(cmd.Body.Discover.Request.EmailAddress, "@")
	if len(email) != 2 {
		// This should NEVER be called because the email is verified in VerifyStructure
		return errors.New("invalid EmailAddress '" + cmd.Body.Discover.Request.EmailAddress + "'")
	}

	for _, domain := range settings.ManagedDomains {
		if email[1] == domain {
			return nil
		}
	}

	return errors.New("the request failed verification") // TODO: Better warning about ManagedDomains
}

func (cmd Request) IsAuthPolicySupport(authPolicy types.AuthPolicy) bool {
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
