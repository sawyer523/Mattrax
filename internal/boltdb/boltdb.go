package boltdb

import (
	"time"

	"github.com/boltdb/bolt"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/certificates"
	"github.com/mattrax/Mattrax/internal/datastore/boltdb"
	"github.com/mattrax/Mattrax/internal/settings"
	"github.com/pkg/errors"
)

// TODO: Use helpers in the package to cut down on duplicate code

// FUTURE: Remove because globals are bad
var globalDB *bolt.DB

// Initialise the database connection and mount the services to the server
func Initialise(server *mattrax.Server) error {
	db, err := bolt.Open(server.Config.DBPath, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		return errors.Wrap(err, "Error initialising Boltdb")
	}
	globalDB = db

	// TODO: Could this be done in parallel to shorten startup time

	if server.UserService, err = NewUserService(db); err != nil {
		return err
	}

	if server.PolicyService, err = NewPolicyService(db); err != nil {
		return err
	}

	settingStore := &boltdb.Store{
		DB:     db,
		Bucket: []byte("settings"),
	}
	settingStore.Init() // TODO: Make to create func not subfunc

	if server.Settings, err = settings.NewService(settingStore); err != nil {
		return err
	}
	if server.Certificates, err = certificates.NewService(settingStore); err != nil {
		return err
	}

	if server.Devices, err = NewDeviceStore(db); err != nil {
		return err
	}

	return nil
}

func Close() error {
	return globalDB.Close()
}
