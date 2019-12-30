package wsettings

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

type Settings struct {
	AuthPolicy          AuthPolicy `graphql:",optional"`
	FederationPortalURL string     `graphql:"federationPortalUrl,optional"` // The URL to handle Federated authentication when its set as the AuthPolicy
	// TODO: Custom ToS URL, Azure AD Creds
}

// Verify checks that the structs fields are valid
func (settings Settings) Verify() error {
	// TODO

	return nil
}
