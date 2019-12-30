package api

import (
	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

// Initialise creates the GraphQL API and attaches its HTTP handler
func Initialise(server *mattrax.Server, r *mux.Router) error {
	// Create Schema
	builder := schemabuilder.NewSchema()

	// Construct Schema Endpoints
	User(server, builder)
	Settings(server, builder)

	// Compile schema
	schema := builder.MustBuild()
	if server.Config.DevelopmentMode {
		introspection.AddIntrospectionToSchema(schema)
	}

	// Mount handler
	r.Handle("/query", graphql.HTTPHandler(schema)).Methods("POST")

	return nil
}
