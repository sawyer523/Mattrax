package settings

import (
	"errors"

	"github.com/imdario/mergo"
	"github.com/rs/zerolog/log"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

// MountAPI attaches the Settings Schema to the GraphQL API
func (s *Service) MountAPI(builder *schemabuilder.Schema) {
	object := builder.Object("Settings", Settings{})
	object.Description = "The dynamic settings that control how your Mattrax server operates."

	var serverStateEnum ServerState
	builder.Enum(serverStateEnum, map[string]ServerState{
		"Installation":       StateInstallation,
		"Normal":             StateNormal,
		"EnrollmentDisabled": StateEnrollmentDisabled,
	})

	query := builder.Query()
	query.FieldFunc("settings", s.Get)

	mutation := builder.Mutation()
	mutation.FieldFunc("updateSettings", func(newSettings Settings) (Settings, error) {
		currentSettings := s.Get()
		if err := mergo.Merge(&newSettings, currentSettings); err != nil {
			log.Error().Err(err).Msg("error merging Settings structs")
			return Settings{}, errors.New("internal server error: failed to merge settings")
		}

		// ServerState is read only through the API.
		if newSettings.ServerState != currentSettings.ServerState {
			return Settings{}, errors.New("the server state is read only")
		}

		err := s.Set(newSettings)
		if err != nil {
			return Settings{}, err
		}
		return newSettings, nil
	})
}
