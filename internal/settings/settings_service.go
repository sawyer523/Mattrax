package settings

import (
	"errors"
	"sync"

	"github.com/imdario/mergo"
	"github.com/mattrax/Mattrax/internal/datastore"
	"github.com/rs/zerolog/log"
)

var settingsKey = []byte("settings")

// Service contains the code for safely (using a Mutex) getting and updating settings.
type Service struct {
	settings Settings
	mutex    *sync.Mutex // Mutex is used to ensures exclusive access to the settings
	store    datastore.Store
}

// Get returns the loaded settings.
func (s *Service) Get() Settings {
	s.mutex.Lock()
	settings := s.settings
	s.mutex.Unlock()

	return settings
}

// Set saves updates settings, if saving fails restore the current values.
func (s *Service) Set(settings Settings) error {
	if err := settings.Verify(); err != nil {
		return err
	}

	currentSettings := s.Get()
	if err := mergo.Merge(&settings, currentSettings); err != nil {
		log.Error().Err(err).Msg("error merging Settings structs")
		return errors.New("internal server error: failed to merge settings")
	}

	s.mutex.Lock()
	previousSettings := s.settings
	s.settings = settings

	if err := s.store.Set(settingsKey, settings); err != nil {
		s.settings = previousSettings
		s.mutex.Unlock()
		log.Error().Err(err).Msg("error saving settings")
		return errors.New("internal error saving settings. values were not changed")
	}
	s.mutex.Unlock()

	return nil
}

// NewService initialises and returns a new SettingsService
func NewService(store datastore.Store) (*Service, error) {
	var settings Settings
	err := store.Get(settingsKey, &settings)
	if err != nil {
		return nil, err
	}

	return &Service{
		settings: settings,
		mutex:    &sync.Mutex{},
		store:    store,
	}, nil
}
