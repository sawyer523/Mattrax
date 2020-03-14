package settings

import (
	"errors"
	"net/url"

	"github.com/mattrax/Mattrax/pkg/types"
)

// Settings contains the Mattrax's dynamic configuration.
// These values can be changed at runtime and they effect how your Mattrax server operates
type Settings struct {
	Tenant               TenantSettings               `json:"tenant"`
	AuthProviderSettings map[string]map[string]string `json:"auth"` // TODO: Make this a better solution that uses typed structs
}

// TenantSettings contains details about the server's owner
// These values may show up on the managed device or during the enrollment process
type TenantSettings struct {
	Name           string `json:"name"`
	SupportEmail   string `json:"support_email"`
	SupportPhone   string `json:"support_phone"`
	SupportWebsite string `json:"support_website"`
	// EnrollmentDisabled bool   `json:"enrollment_disabled"`
}

// Verify checks the Settings are valid. This is done prior to saving updated settings.
func (settings Settings) Verify() error {
	if settings.Tenant.Name != "" && !types.GenericStringRegex.MatchString(settings.Tenant.Name) {
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
