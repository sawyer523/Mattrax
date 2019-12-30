package boltdb

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/pkg/errors"
)

// settingsBucket stores the name of the boltdb bucket the settings are stored in
var settingsBucket = []byte("settings")

// settingsKey is the key within the database that the settings are stored to
var settingsKey = []byte("settings")

// SettingsStore saves and loads the servers settings
type SettingsStore struct {
	db *bolt.DB
}

// Retrieve gets the settings from the database
func (st SettingsStore) Retrieve() (types.Settings, error) {
	var settings types.Settings
	err := st.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(settingsBucket)
		if bucket == nil {
			return errors.New("error in SettingsService.Get: settings bucket does not exist")
		}

		settingsRaw := bucket.Get(settingsKey)
		if settingsRaw == nil {
			// Blank settings returned
			return nil
		}

		err := gob.NewDecoder(bytes.NewBuffer(settingsRaw)).Decode(&settings)
		return err
	})

	return settings, err
}

// Save stores new settings to the database
func (st SettingsStore) Save(settings types.Settings) error {
	// Encode Settings
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(settings); err != nil {
		return errors.Wrap(err, "error in SettingsService.Update: problem to encoding settings struct")
	}
	settingsRaw := buf.Bytes()

	// Store to DB
	err := st.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(settingsBucket)
		if bucket == nil {
			return errors.New("error in SettingsService.Update: settings bucket does not exist")
		}

		err := bucket.Put(settingsKey, settingsRaw)
		return err
	})

	return err
}

// NewSettingsStore creates and initialises a new SettingsService from a DB connection
func NewSettingsStore(db *bolt.DB) (types.SettingsStore, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(settingsBucket)
		return err
	})

	return SettingsStore{
		db,
	}, err
}
