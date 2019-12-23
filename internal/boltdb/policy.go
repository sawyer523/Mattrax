package boltdb

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/pkg/errors"
)

// policiesBucket stores the name of the boltdb bucket the polcies are stored in
var policiesBucket = []byte("policies")

// PolicyService contains the implemented functionality for policies
type PolicyService struct {
	db *bolt.DB
}

// GetAll returns all policies
func (ps PolicyService) GetAll() ([]types.Policy, error) {
	var policies []types.Policy
	err := ps.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(policiesBucket)
		if bucket == nil {
			return errors.New("error in PolicyService.GetAll: policies bucket does not exist")
		}

		c := bucket.Cursor()
		for key, policyRaw := c.First(); key != nil; key, policyRaw = c.Next() {
			var policy types.Policy
			err := gob.NewDecoder(bytes.NewBuffer(policyRaw)).Decode(&policy)
			if err != nil {
				return errors.Wrap(err, "error in PolicyService.GetAll: problem to decoding the policy struct")
			}

			policies = append(policies, policy)
		}

		return nil
	})

	return policies, err
}

// Get returns a policy by its uuid
func (ps PolicyService) Get(uuid types.PolicyUUID) (types.Policy, error) {
	var policy types.Policy
	err := ps.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(policiesBucket)
		if bucket == nil {
			return errors.New("error in PolicyService.Get: policies bucket does not exist")
		}

		policyRaw := bucket.Get(uuid)
		if policyRaw == nil {
			fmt.Println("NULL")
			return nil // TODO: Cpstom Exported Error
		}

		err := gob.NewDecoder(bytes.NewBuffer(policyRaw)).Decode(&policy)

		return err
	})

	return policy, err
}

// CreateOrEdit creates or edits an existing policy if one exists
func (ps PolicyService) CreateOrEdit(uuid types.PolicyUUID, policy types.Policy) error {
	// Encode policy
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(policy); err != nil {
		return errors.Wrap(err, "error in PolicyService.CreateOrEdit: problem to encoding policy struct")
	}
	policyRaw := buf.Bytes()

	// Store to DB
	err := ps.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(policiesBucket)
		if bucket == nil {
			return errors.New("error in PolicyService.CreateOrEdit: policies bucket does not exist")
		}

		err := bucket.Put(uuid, policyRaw)
		return err
	})

	return err
}

// GenerateUUID returns a new PolicyUUID
func (ps PolicyService) GenerateUUID() (types.PolicyUUID, error) {
	// TODO: Make Do Stuffz
	return types.PolicyUUID([]byte{}), nil
}

// NewPolicyService creates and initialises a new PolicyService from a DB connection
func NewPolicyService(db *bolt.DB) (PolicyService, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(policiesBucket)
		return err
	})

	return PolicyService{
		db,
	}, err
}
