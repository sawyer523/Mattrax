package boltdb

import (
	"crypto/x509/pkix"
	"time"

	"github.com/boltdb/bolt"
	mattrax "github.com/mattrax/Mattrax/internal"
	certificateservice "github.com/mattrax/Mattrax/internal/certificates"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/pkg/errors"
)

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

	if server.SettingsStore, err = NewSettingsStore(db); err != nil {
		return err
	}

	if server.Settings, err = server.SettingsStore.Retrieve(); err != nil {
		return err
	}

	identityCertConfig := types.IdentityCertificateConfig{
		KeyLength: 4096,
		Subject: pkix.Name{
			// TODO: Configurable
			Country:            []string{"US"},
			Organization:       []string{"groob-io"},
			OrganizationalUnit: []string{"SCEP CA"},
		},
	}

	if server.CertificateStore, err = NewCertificateStore(db); err != nil {
		return err
	}

	if server.Certificates, err = certificateservice.Initialise(server.CertificateStore, identityCertConfig); err != nil {
		return err
	}

	return nil
}

func Close() error {
	return globalDB.Close()
}
