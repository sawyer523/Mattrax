package settings

import (
	"errors"
	"github.com/rs/zerolog/log"
	"sync"
)

// Service contains the code for safely (using a Mutex) getting and updating settings.
type Service struct {
	settings Settings
	mutex    *sync.Mutex // Mutex is used to ensures exclusive access to the settings
	store    Store
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

	s.mutex.Lock()
	previousSettings := s.settings
	s.settings = settings

	if err := s.store.Save(settings); err != nil {
		s.settings = previousSettings
		s.mutex.Unlock()
		log.Error().Err(err).Msg("error saving settings")
		return errors.New("internal error saving settings. values were not changed")

	}
	s.mutex.Unlock()

	return nil
}

// Store is a place where settings are stored.
type Store interface {
	Save(Settings) error
	Retrieve() (Settings, error)
}

// NewService initialises and returns a new SettingsService
func NewService(store Store) (*Service, error) {
	settings, err := store.Retrieve()
	if err != nil {
		return nil, err
	}

	return &Service{
		settings: settings,
		mutex:    &sync.Mutex{},
		store:    store,
	}, nil
}
