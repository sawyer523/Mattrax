package types

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
	AuthPolicy AuthPolicy `graphql:",optional"`
	FederationPortalURL string `graphql:"federationPortalUrl,optional"` // The URL to handle Federated authentication when its set as the AuthPolicy
}

// Settings contains the global settings for your Mattrax server
type Settings struct {
	TenantName   string `graphql:",optional"`
	ManagedDomains []string `graphql:",optional"`
	Windows WindowsSettings `graphql:",optional"`
}

// SettingsService contains the implemented functionality for managing settings
type SettingsService interface {
	Get() (Settings, error)
	Update(settings Settings) error
}
