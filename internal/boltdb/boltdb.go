package boltdb

import (
	"time"

	"github.com/boltdb/bolt"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/certificates"
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

	if settingsStore, err := NewSettingsStore(db); err != nil {
		return err
	} else if server.Settings, err = settings.NewService(settingsStore); err != nil {
		return err
	}

	if certificateStore, err := NewCertificateStore(db); err != nil {
		return err
	} else if server.Certificates, err = certificates.NewService(certificateStore); err != nil {
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
