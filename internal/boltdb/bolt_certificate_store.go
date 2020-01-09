package boltdb

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
	"github.com/mattrax/Mattrax/internal/certificates"
	"github.com/pkg/errors"
)

// TODO: Make smaller + Redo error messages. Log them to console. // Same with other BoltDB stores.

// Where in the DB the certificates struct is stored
var certificatesBucket = []byte("certificates")
var certificatesKey = []byte("certificates")

// CertificateStore saves and loads the servers certificates
type CertificateStore struct {
	db *bolt.DB
}

// Retrieve gets the certificates from the database
func (cs CertificateStore) Retrieve() (certificates.Certificates, error) {
	var certificates certificates.Certificates
	err := cs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(certificatesBucket)
		if bucket == nil {
			return errors.New("error in CertificateStore.Retrieve: certificates bucket does not exist")
		}

		certificatesRaw := bucket.Get(certificatesKey)
		if certificatesRaw == nil {
			// Blank certificates returned
			return nil
		}

		err := gob.NewDecoder(bytes.NewBuffer(certificatesRaw)).Decode(&certificates)
		return err
	})

	return certificates, err
}

// Save stores certificates to the database
func (cs CertificateStore) Save(certificates certificates.Certificates) error {
	// Encode Certificates
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(certificates); err != nil {
		return errors.Wrap(err, "error in CertificatesStore.Save: problem to encoding certificates struct")
	}
	certificatesRaw := buf.Bytes()

	// Store to DB
	err := cs.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(certificatesBucket)
		if bucket == nil {
			return errors.New("error in CertificatesStore.Update: certificates bucket does not exist")
		}

		err := bucket.Put(certificatesKey, certificatesRaw)
		return err
	})

	return err
}

// NewCertificateStore creates and initialises a new CertificateStore from a DB connection
func NewCertificateStore(db *bolt.DB) (certificates.Store, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(certificatesBucket)
		return err
	})

	return CertificateStore{
		db,
	}, err
}
