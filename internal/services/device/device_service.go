package device

import (
	"errors"

	"github.com/imdario/mergo"
	"github.com/mattrax/Mattrax/internal/datastore"
	"github.com/rs/zerolog/log"
)

// Service exposes the Devices from the underlying datastore
type Service struct {
	store datastore.Zone // The underlying datastore zone to save the devices into
}

// GetAll returns all of the enrolled devices
func (s *Service) GetAll() (map[string]*Device, error) {
	devices, err := s.store.GetAll(Device{})
	return devices.Interface().(map[string]*Device), err
}

// Get returns an enrolled device from its uuid
func (s *Service) Get(uuid string) (Device, error) {
	var device Device
	return device, s.store.Get(uuid, &device)
}

// SaveSS saves a device to the datastore.
// The SS means server side and this is because no
// restrictions are applied to the data being entered
func (s *Service) SaveSS(uuid string, device Device) error {
	return s.store.Set(uuid, device)
}

// Save updates a device
func (s *Service) Save(uuid string, device Device) error {
	if err := device.Verify(); err != nil {
		return err
	}

	currentDevice, err := s.Get(uuid)
	if err != nil {
		return err
	}
	if err := mergo.Merge(&device, currentDevice); err != nil {
		log.Error().Err(err).Msg("error merging Settings structs")
		return errors.New("internal server error: failed to merge settings")
	}

	if err := s.store.Set(uuid, device); err != nil {
		log.Error().Err(err).Msg("error saving device")
		return errors.New("internal error saving device")
	}

	return nil
}

// Save updates a device
func (s *Service) Unenroll(uuid string) error {
	_, err := s.Get(uuid)
	if err != nil {
		return err
	}
	// TODO: Send signals/etc to device to remove it them once it reports sucess remove from db!

	return s.store.Delete(uuid)
}

// NewService initialises and returns a new DeviceService
func NewService(store datastore.Zone) (*Service, error) {
	return &Service{
		store: store,
	}, nil
}
