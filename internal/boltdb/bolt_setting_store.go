package boltdb

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
	"github.com/mattrax/Mattrax/internal/settings"
	"github.com/pkg/errors"
)

// TODO: Make smaller + Redo error messages. Log them to console.

// Where in the DB the settings struct is stored
var settingsBucket = []byte("settings")
var settingsKey = []byte("settings")

// SettingsStore saves and loads the servers settings
type SettingsStore struct {
	db *bolt.DB
}

// Retrieve gets the settings from the database
func (ss SettingsStore) Retrieve() (settings.Settings, error) {
	var settings settings.Settings
	err := ss.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(settingsBucket)
		if bucket == nil {
			return errors.New("error in SettingsStore.Retrieve: settings bucket does not exist")
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
func (ss SettingsStore) Save(settings settings.Settings) error {
	// Encode Settings
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(settings); err != nil {
		return errors.Wrap(err, "error in SettingsStore.Save: problem to encoding settings struct")
	}
	settingsRaw := buf.Bytes()

	// Store to DB
	err := ss.db.Update(func(tx *bolt.Tx) error {
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
func NewSettingsStore(db *bolt.DB) (settings.Store, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(settingsBucket)
		return err
	})

	return SettingsStore{
		db,
	}, err
}
