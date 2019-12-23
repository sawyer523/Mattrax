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

	mutation := builder.Mutation()
	mutation.FieldFunc("updateSettings", func(settings types.Settings) types.Settings {
		err := server.SettingsService.Update(settings)
		if err != nil {
			panic(err) // TODO: Hande
		}
		return settings
	})
}
