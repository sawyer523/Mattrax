package wsettings

type Settings struct {
	DeploymentType      DeploymentType `graphql:",optional"`
	FederationPortalURL string         `graphql:"federationPortalUrl,optional"` // Federation Portal URL. Leave blank for Mattrax default provider.
	// TODO: Custom ToS URL, Azure AD Creds
}

// Verify checks that the structs fields are valid
func (settings Settings) Verify() error {
	// TODO

	return nil
}
