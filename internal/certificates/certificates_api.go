package certificates

import (
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

// MountAPI attaches the Certificates Schema to the GraphQL API
func (s *Service) MountAPI(builder *schemabuilder.Schema) {
	object := builder.Object("Certificates", Certificates{})
	object.Description = "The details for certificates used to keep Mattrax running and secure."

	identityObject := builder.Object("Identity", Identity{})
	identityObject.FieldFunc("subject", func() string { return s.certificates.Identity.Subject.String() })

	query := builder.Query()
	query.FieldFunc("certificates", s.Get)
}
