package api

import (
	"errors"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	wsettings "github.com/mattrax/Mattrax/mdm/windows/settings"
	"github.com/rs/zerolog/log"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

// Settings contains the API endpoints and types for Settings
func Settings(server *mattrax.Server, builder *schemabuilder.Schema) {
	settingsObject := builder.Object("Settings", types.Settings{})
	settingsObject.Description = "The Dynamic Settings For the Mattrax Server."

	var enumField wsettings.AuthPolicy
	builder.Enum(enumField, map[string]wsettings.AuthPolicy{
		"OnPremise":   wsettings.AuthPolicyOnPremise,
		"Federated":   wsettings.AuthPolicyFederated,
		"Certificate": wsettings.AuthPolicyCertificate,
	})

	query := builder.Query()
	query.FieldFunc("settings", func() types.Settings { return server.Settings })

	mutation := builder.Mutation()
	mutation.FieldFunc("updateSettings", func(settings types.Settings) (types.Settings, error) {
		if err := settings.Verify(); err != nil {
			return types.Settings{}, err
		}

		if err := server.SettingsStore.Save(settings); err != nil {
			log.Error().Err(err).Msg("Error saving settings!")
			return types.Settings{}, errors.New("Error saving the settings")
		}

		server.Settings = settings

		return settings, nil
	})
}
