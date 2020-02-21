package settings

import (
	"errors"
	"net/url"
	"regexp"
)

// Settings contains the Mattrax's dynamic configuration.
// These values can be changed at runtime although it is recommended some of them never change.
type Settings struct {
	Tenant TenantSettings `yaml:"tenant"`
}

// TenantSettings contains details about the server's owner
// Some of these settings show up on the device to tell a end user where to contact for help.
type TenantSettings struct {
	Name               string `yaml:"name"`
	SupportEmail       string `yaml:"support_email"`
	SupportPhone       string `yaml:"support_phone"`
	SupportWebsite     string `yaml:"support_website"`
	EnrollmentDisabled bool   `yaml:"enrollment_disabled"`
}

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

	return nil
}
