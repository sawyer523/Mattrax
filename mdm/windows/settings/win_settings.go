package winsettings

// Settings contains Windows MDM specific settings
type Settings struct {
	CustomFederationPortal string // Replace the Mattrax Federation portal. Note: Currently not fully implemented
	// TODO: CustomTermsOfServicePortal string // Replace the Mattrax AzureAD Terms of Service portal.
}
