package types

// Settings contains the global settings for your Mattrax server
type Settings struct {
	TenantName   string `graphql:",optional"`
	ManagedDomains []string `graphql:",optional"`
}

// SettingsService contains the implemented functionality for managing settings
type SettingsService interface {
	Get() (Settings, error)
	Update(settings Settings) error
}
