package datastore

type Store interface {
	Init() error
	Set(key []byte, value interface{}) error
	Get(key []byte, value interface{}) error
}
