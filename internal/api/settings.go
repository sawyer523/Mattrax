package api

import (
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

func settings(server mattrax.Server, builder *schemabuilder.Schema) {
	settingsObject := builder.Object("Settings", types.Settings{})
	settingsObject.Description = "A Mattrax Server's Settings."

	query := builder.Query()
	query.FieldFunc("settings", server.SettingsService.Get)

	var enumField types.AuthPolicy
	builder.Enum(enumField, map[string]types.AuthPolicy{
		"OnPremise":   types.AuthPolicyOnPremise,
		"Federated":   types.AuthPolicyFederated,
		"Certificate": types.AuthPolicyCertificate,
	})

	mutation := builder.Mutation()
	mutation.FieldFunc("updateSettings", func(settings types.Settings) (types.Settings, error) {
		if err := settings.Verify(); err != nil {
			return types.Settings{}, err
		}

		if err := server.SettingsService.Update(settings); err != nil {
			panic(err) // TODO: Handle
		}

		return settings, nil
	})
}
