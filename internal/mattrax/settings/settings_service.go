package settings

import (
	"errors"
	"sync"

	"github.com/imdario/mergo"
	"github.com/mattrax/Mattrax/internal/datastore"
	"github.com/rs/zerolog/log"
)

// Service exposes the Settings from the underlying datastore
type Service struct {
	settings Settings       // Cached settings
	mutex    *sync.Mutex    // Mutex is used to prevent multiple go routines using settings at the same time causing a race condition
	store    datastore.Zone // The underlying datastore zone to save the settings into
}

// Get returns current settings
func (s *Service) Get() Settings {
	s.mutex.Lock()
	settings := s.settings
	s.mutex.Unlock()
	return settings
}

// Set updates the settings, restoring to previous in the event of failure
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

	if err := s.store.Set("settings", settings); err != nil {
		s.settings = previousSettings
		s.mutex.Unlock()
		log.Error().Err(err).Msg("error saving settings")
		return errors.New("internal error saving settings. values were not changed")
	}
	s.mutex.Unlock()

	return nil
}

// NewService initialises and returns a new SettingsService
func NewService(store datastore.Zone) (*Service, error) {
	settings := Settings{}
	if err := store.Get("settings", &settings); err != nil && err != datastore.ErrNotFound {
		return nil, err
	}

	return &Service{
		settings: settings,
		mutex:    &sync.Mutex{},
		store:    store,
	}, nil
}
