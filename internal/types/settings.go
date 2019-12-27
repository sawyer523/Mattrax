package types

import (
	"net/url"
	"regexp"

	"github.com/pkg/errors"
)

// AuthPolicy is the method Windows uses to authentication the client
type AuthPolicy int

const (
	// AuthPolicyOnPremise is the OnPremise Windows AuthPolicy
	AuthPolicyOnPremise AuthPolicy = iota
	// AuthPolicyFederated is the Federated Windows AuthPolicy
	AuthPolicyFederated
	// AuthPolicyCertificate is the Certificate Windows AuthPolicy
	AuthPolicyCertificate
)

// WindowsSettings contains Windows MDM specific settings
// TODO: Move to /mdm/windows/ folder
type WindowsSettings struct {
	AuthPolicy          AuthPolicy `graphql:",optional"`
	FederationPortalURL string     `graphql:"federationPortalUrl,optional"` // The URL to handle Federated authentication when its set as the AuthPolicy
}

// Settings contains the global settings for your Mattrax server
type Settings struct {
	TenantName     string          `graphql:",optional"`
	ManagedDomains []string        `graphql:",optional"`
	Windows        WindowsSettings `graphql:",optional"`
}

// tenantName is a regex used to verify a tenant name is valid
var tenantName = regexp.MustCompile(`^[a-zA-Z0-9- '"]+$`)

func (settings Settings) Verify() error {
	if settings.TenantName != "" && !tenantName.MatchString(settings.TenantName) {
		return errors.New("invalid settings: invalid TenantName '" + settings.TenantName + "'")
	}

	// TODO: Verify ManagedDomains

	if settings.Windows.AuthPolicy > 3 {
		return errors.New("invalid settings: invalid AuthPolicy")
	}

	if settings.Windows.FederationPortalURL != "" {
		if federationPortalURL, err := url.ParseRequestURI(settings.Windows.FederationPortalURL); err != nil {
			return errors.New("invalid settings: invalid FederationPortalURL '" + settings.Windows.FederationPortalURL + "'")
		} else if federationPortalURL.Scheme != "https" {
			return errors.New("invalid settings: invalid FederationPortalURL scheme '" + federationPortalURL.Scheme + "', It must be 'https'")
		}
	}

	return nil
}

// SettingsService contains the implemented functionality for managing settings
type SettingsService interface {
	Get() (Settings, error)
	Update(settings Settings) error
}
