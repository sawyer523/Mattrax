package settings

import (
	"errors"
	winsettings "github.com/mattrax/Mattrax/mdm/windows/settings"
	"net/url"
	"regexp"
)

// Settings contains the Mattrax's dynamic configuration.
// These values are changed by the Tenant at runtime.
type Settings struct {
	Tenant      TenantSettings `graphql:",optional"`
	ServerState ServerState    `graphql:",optional"`
	// TODO UserSource string

	Windows winsettings.Settings `graphql:",optional"`
}

// TenantSettings contains settings that are specific to the servers tenant
type TenantSettings struct {
	Name           string `graphql:",optional"`
	SupportPhone   string `graphql:",optional"`
	SupportEmail   string `graphql:",optional"`
	SupportWebsite string `graphql:",optional"`
}

// ServerState says the state of the server
type ServerState int

const (
	// StateInstallation is the ServerState while the MDM is still being configurated upon instllation
	StateInstallation ServerState = iota
	// StateNormal is the ServerState for the MDM server functioning normally
	StateNormal
	// StateEnrollmentDisabled is the ServerState for when functioning normally but device enrollment is disabled
	StateEnrollmentDisabled
)

// genericStringRegex is a regex used to verify a simple string
var genericStringRegex = regexp.MustCompile(`^[a-zA-Z0-9- '"]+$`)

// Verify checks the Settings are valid. This is done prior to saving updated settings.
func (settings Settings) Verify() error {
	if settings.Tenant.Name != "" && !genericStringRegex.MatchString(settings.Tenant.Name) {
		return errors.New("invalid settings: tenant name contains invalid characters")
	}

	// TODO: Verify SupportPhone + SupportEmail

	if settings.Tenant.SupportWebsite != "" {
		if _, err := url.ParseRequestURI(settings.Tenant.SupportWebsite); err != nil {
			return errors.New("invalid settings: tenant name contains invalid characters")
		}
	}

	// TODO

	return nil
}
