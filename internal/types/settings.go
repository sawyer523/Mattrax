package types

import (
	"errors"
	"regexp"

	wsettings "github.com/mattrax/Mattrax/mdm/windows/settings"
)

// Settings holds the dynamic server config
// This can be changed via the API at runtime.
type Settings struct {
	TenantName        string             `graphql:",optional"`
	ManagedDomains    []string           `graphql:",optional"`
	EnrollmentEnabled bool               `graphql:",optional"`
	Windows           wsettings.Settings `graphql:",optional"`
}

// Regex's are used to verify the users input
var tenantNameRegex = regexp.MustCompile(`^[a-zA-Z0-9- '"]+$`)

// IsDNSNameRegex is used to verify if a string is a DNS Name (Domain Name)
var IsDNSNameRegex = regexp.MustCompile(`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`)

// Verify checks that the structs fields are valid
func (settings Settings) Verify() error {
	if settings.TenantName != "" && !tenantNameRegex.MatchString(settings.TenantName) {
		return errors.New("invalid settings: invalid TenantName '" + settings.TenantName + "'")
	}

	for _, domain := range settings.ManagedDomains {
		if !IsDNSNameRegex.MatchString(domain) {
			return errors.New("invalid settings: invalid ManagedDomain '" + domain + "'")
		}
	}

	err := settings.Windows.Verify()
	return err
}

// SettingsStore is a storage mechanise capable of permanently storing settings
type SettingsStore interface {
	Retrieve() (Settings, error)
	Save(settings Settings) error
}
