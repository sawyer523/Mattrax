package api

import (
	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

// InitAPI initialises the GraphQL API
func InitAPI(server mattrax.Server, r *mux.Router) {
	// Build Schema
	builder := schemabuilder.NewSchema()
	user(server, builder)
	settings(server, builder)

	// Compile Schema
	schema := builder.MustBuild()
	if server.Config.DevelopmentMode {
		introspection.AddIntrospectionToSchema(schema)
	}

	// Serve GraphQL HTTP Endpoint
	r.Handle("/query", graphql.HTTPHandler(schema)).Methods("POST")
}
