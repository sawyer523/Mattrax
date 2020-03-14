package datastore

import (
	"errors"
	"reflect"
)

// ErrNotFound is returned from Get() when the item is not found
var ErrNotFound = errors.New("the key you requested does not exist in the current zone")

// Store represents a backend for storing Mattrax's data.
// This struct is contained to use in the "mattrax" package and
// individual zones are parsed to services as they require them.
type Store interface {
	Init(customConnStr string) error // customConnStr is the part after the ":"
	Close() error
	Zone(name string) (Zone, error) // Zone is equivalent to a boltdb bucket or table
}

// Zone holds a single type of data
type Zone interface {
	Set(key string, value interface{}) error
	Get(key string, model interface{}) error
	GetAll(model interface{}) (reflect.Value, error)
	Delete(key string) error
}
