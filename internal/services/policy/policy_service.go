package policy

import (
	"errors"

	"github.com/imdario/mergo"
	"github.com/mattrax/Mattrax/internal/datastore"
	"github.com/rs/zerolog/log"
)

// Service exposes the Policys from the underlying datastore
type Service struct {
	store datastore.Zone // The underlying datastore zone to save the policys into
}

// GetAll returns all of the enrolled policys
func (s *Service) GetAll() (map[string]*Policy, error) {
	policys, err := s.store.GetAll(Policy{})
	return policys.Interface().(map[string]*Policy), err
}

// Get returns an enrolled policy from its uuid
func (s *Service) Get(uuid string) (Policy, error) {
	var policy Policy
	return policy, s.store.Get(uuid, &policy)
}

// SaveSS saves a policy to the datastore.
// The SS means server side and this is because no
// restrictions are applied to the data being entered
func (s *Service) SaveSS(uuid string, policy Policy) error {
	return s.store.Set(uuid, policy)
}

// Save updates a policy
func (s *Service) Save(uuid string, policy Policy) error {
	if err := policy.Verify(); err != nil {
		return err
	}

	currentPolicy, err := s.Get(uuid)
	if err == datastore.ErrNotFound {
		currentPolicy = Policy{}
	} else if err != nil {
		return err
	}
	if err := mergo.Merge(&policy, currentPolicy); err != nil {
		log.Error().Err(err).Msg("error merging Settings structs")
		return errors.New("internal server error: failed to merge settings")
	}

	if err := s.store.Set(uuid, policy); err != nil {
		log.Error().Err(err).Msg("error saving policy")
		return errors.New("internal error saving policy")
	}

	return nil
}

// Delete removes a policy
func (s *Service) Delete(uuid string) error {
	_, err := s.Get(uuid)
	if err != nil {
		return err
	}

	return s.store.Delete(uuid)
}

// NewService initialises and returns a new PolicyService
func NewService(store datastore.Zone) (*Service, error) {
	return &Service{
		store: store,
	}, nil
}
