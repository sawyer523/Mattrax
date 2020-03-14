package boltdb

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mattrax/Mattrax/internal/datastore"
	"github.com/pkg/errors"
)

type Store struct {
	db *bolt.DB
}

func (s *Store) Init(customConnURL string) error {
	db, err := bolt.Open(customConnURL, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return errors.Wrap(err, "Error initialising Boltdb")
	}
	s.db = db

	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Zone(name string) (datastore.Zone, error) {
	zone := &Zone{
		db:     s.db,
		Bucket: []byte(name),
	}

	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(zone.Bucket)
		return err
	})
	if err != nil {
		return nil, errors.Wrap(err, "boltdb: error initialising zone bucket")
	}

	return zone, nil
}

type Zone struct {
	db     *bolt.DB
	Bucket []byte
}

func (z *Zone) Set(key string, value interface{}) error {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(value); err != nil {
		return errors.Wrap(err, "boltdb: unable to encode value")
	}
	raw := buf.Bytes()

	err := z.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(z.Bucket)
		if bucket == nil {
			return errors.New("boltdb: bucket does not exist")
		}

		err := bucket.Put([]byte(key), raw)
		return err
	})

	return err
}

func (z *Zone) Get(key string, model interface{}) error {
	err := z.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(z.Bucket)
		if bucket == nil {
			return errors.New("boltdb: bucket does not exist")
		}

		raw := bucket.Get([]byte(key))
		if raw == nil {
			return datastore.ErrNotFound
		}

		err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(model)
		return err
	})

	return err
}

func (z *Zone) GetAll(model interface{}) (reflect.Value, error) {
	out := reflect.MakeMap(reflect.MapOf(reflect.TypeOf(""), reflect.PtrTo(reflect.TypeOf(model))))
	err := z.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(z.Bucket)
		if bucket == nil {
			return errors.New("boltdb: bucket does not exist")
		}

		return bucket.ForEach(func(uuid, raw []byte) error {
			var item = reflect.New(reflect.TypeOf(model)).Interface()
			err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(item)
			if err != nil {
				return err
			}

			out.SetMapIndex(reflect.ValueOf(string(uuid)), reflect.ValueOf(item))
			return nil
		})
	})
	return out, err
}

func (z *Zone) Delete(key string) error {
	err := z.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(z.Bucket)
		if bucket == nil {
			return errors.New("boltdb: bucket does not exist")
		}

		return bucket.Delete([]byte(key))
	})
	return err
}
