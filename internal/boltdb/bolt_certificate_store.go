package boltdb

import (
	"crypto/rsa"
	"crypto/x509"

	"github.com/boltdb/bolt"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// certificatesBucket stores the name of the boltdb bucket that the certificates are stored in
var certificatesBucket = []byte("certificates")

// The bucket key used when storing
var identityCertificateKey = []byte("identityCertificate")
var identityPrivateKeyKey = []byte("identityPrivateKey")

// CertificateStore saves and loads the servers certificates
type CertificateStore struct {
	db *bolt.DB
}

// RetrieveIdentity gets the identity certficate & private key and their raw values from the datastore
func (cs CertificateStore) RetrieveIdentity() ([]byte, *x509.Certificate, []byte, *rsa.PrivateKey, error) {
	var certificateDer []byte
	var privateKeyDer []byte
	if err := cs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(certificatesBucket)
		if bucket == nil {
			log.Error().Err(errors.New("Certificates bucket doesn't exist")).Msg("Error retrieving identity certificates!")
			return errors.New("internal server error: please get your administrator to check the logs")
		}

		certificateDer = bucket.Get(identityCertificateKey)
		privateKeyDer = bucket.Get(identityPrivateKeyKey)

		return nil
	}); err != nil {
		return nil, nil, nil, nil, err
	}

	if certificateDer == nil || privateKeyDer == nil {
		return nil, nil, nil, nil, types.ErrIdentityNotFound
	}

	cert, err := x509.ParseCertificate(certificateDer)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "error unable to decode the identity certificate")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyDer)
	if err != nil {
		return nil, nil, nil, nil, errors.Wrap(err, "error unable to decode the identity private key")
	}

	return certificateDer, cert, privateKeyDer, privateKey, err
}

// SaveIdentity saves the identity certficate & private key to the datastore
func (cs CertificateStore) SaveIdentity(cert []byte, key *rsa.PrivateKey) error {
	if err := cs.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(certificatesBucket)
		if bucket == nil {
			return errors.New("error certificates bucket does not exist")
		}

		if err := bucket.Put(identityCertificateKey, cert); err != nil {
			return err
		}
		if err := bucket.Put(identityPrivateKeyKey, x509.MarshalPKCS1PrivateKey(key)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// NewCertificateStore creates and initialises a new NewCertificateStore from a DB connection
func NewCertificateStore(db *bolt.DB) (CertificateStore, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(certificatesBucket)
		if err != nil {
			return errors.Wrap(err, "error NewCertificateStore: error creating bucket")
		}

		return nil
	})

	return CertificateStore{
		db,
	}, err
}
