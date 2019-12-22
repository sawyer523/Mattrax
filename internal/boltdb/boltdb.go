package boltdb

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// Init initialises the database connection
func Init() (*bolt.DB, error) {
	db, err := bolt.Open("./mattrax.db", 0600, &bolt.Options{Timeout: 5 * time.Second})
	return db, errors.Wrap(err, "Error initialising Boltdb")
}
