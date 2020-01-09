package api

import (
	"context"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/devices"
	"github.com/mattrax/Mattrax/internal/settings"
	"github.com/samsarahq/thunder/graphql"
	"github.com/samsarahq/thunder/graphql/introspection"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
	"github.com/samsarahq/thunder/reactive"
)

// Initialise creates the GraphQL API and attaches its HTTP handler
func Initialise(server *mattrax.Server, r *mux.Router) error {
	// Create Schema
	builder := schemabuilder.NewSchema()

	// Construct Schema Endpoints
	MattraxAPI(server, builder)
	server.Settings.MountAPI(builder)
	server.Certificates.MountAPI(builder)
	devices.MountAPI(server.Devices, builder)

	// TODO: Move below to somewhere else
	type finishInstallation struct {
		Underscore bool `graphql:"_"`
	}

	mutation := builder.Mutation()
	mutation.FieldFunc("finishInstallation", func(req finishInstallation) (finishInstallation, error) {
		currentSettings := server.Settings.Get()
		fmt.Println(currentSettings)
		if currentSettings.ServerState != settings.StateInstallation {
			return finishInstallation{false}, errors.New("failed to finish instllation: server is not in installation mode")
		}

		if currentSettings.Tenant.Name == "" {
			return finishInstallation{false}, errors.New("Failed to finish installation: Tenant name must be set to finish enrollment")
		}

		if err := server.Certificates.GenerateIdentity(pkix.Name{
			CommonName: currentSettings.Tenant.Name + " Identity",
		}); err != nil {
			return finishInstallation{false}, err
		}

		if err := server.Settings.Set(settings.Settings{
			ServerState: settings.StateNormal,
		}); err != nil {
			return finishInstallation{false}, err
		}

		return finishInstallation{true}, nil
	})

	// Compile schema
	schema := builder.MustBuild()
	if server.Config.DevelopmentMode {
		introspection.AddIntrospectionToSchema(schema)
	}

	// Mount handler
	r.Handle("/query", graphql.HTTPHandler(schema)).Methods("POST")

	return nil
}

// MattraxAPI exposes values that are directly stored on the mattrax.Server to the API
func MattraxAPI(server *mattrax.Server, builder *schemabuilder.Schema) {
	object := builder.Object("Config", mattrax.Config{})
	object.Description = "The static server config. This is set via command line flags and is read only."

	query := builder.Query()
	query.FieldFunc("version", func(ctx context.Context) string {
		reactive.InvalidateAfter(ctx, 24*60*time.Second)
		return server.Version
	})
	query.FieldFunc("config", func(ctx context.Context) mattrax.Config {
		reactive.InvalidateAfter(ctx, 24*60*time.Second)
		return server.Config
	})
}
