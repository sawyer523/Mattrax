package boltdb

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type Store struct {
	DB     *bolt.DB
	Bucket []byte
}

func (s *Store) Init() error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(s.Bucket)
		return err
	})
}

func (s *Store) Set(key []byte, value interface{}) error {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(value); err != nil {
		return errors.Wrap(err, "error: boltdb: unable to encode value")
	}
	raw := buf.Bytes()

	err := s.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.Bucket)
		if bucket == nil {
			return errors.New("error: boltdb: bucket does not exist")
		}

		err := bucket.Put(key, raw)
		return err
	})

	return err
}

func (s *Store) Get(key []byte, model interface{}) error {
	err := s.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(s.Bucket)
		if bucket == nil {
			return errors.New("error: boltdb: bucket does not exist")
		}

		raw := bucket.Get(key)
		if raw == nil {
			return nil
		}

		err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(model)
		return err
	})

	return err
}
