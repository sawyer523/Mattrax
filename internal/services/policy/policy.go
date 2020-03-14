package policy

// Policy represents an MDM enrolled policy
type Policy struct {
	DisplayName string `json:"name"`
}

// Verify checks the Policy is valid. This is done prior to saving updated settings.
func (policy Policy) Verify() error {
	// 	if settings.Tenant.Name != "" && !types.GenericStringRegex.MatchString(settings.Tenant.Name) {
	// 		return errors.New("invalid settings: tenant name contains invalid characters")
	// 	}

	// 	// TODO: Verify SupportPhone + SupportEmail

	// 	if settings.Tenant.SupportWebsite != "" {
	// 		if _, err := url.ParseRequestURI(settings.Tenant.SupportWebsite); err != nil {
	// 			return errors.New("invalid settings: tenant name contains invalid characters")
	// 		}
	// 	}

	return nil
}
