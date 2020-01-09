package boltdb

import (
	"bytes"
	"encoding/gob"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/mattrax/Mattrax/internal/devices"
	"github.com/pkg/errors"
)

// devicesBucket stores the name of the boltdb bucket the users are stored in
var devicesBucket = []byte("devices")

// DeviceStore saves and loads devices
type DeviceStore struct {
	db *bolt.DB
}

// GetAll returns all devices
func (ds DeviceStore) GetAll() ([]devices.Device, error) {
	var devicesList []devices.Device
	err := ds.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(devicesBucket)
		if bucket == nil {
			return errors.New("error in DeviceStore.GetAll: devices bucket does not exist")
		}

		c := bucket.Cursor()
		for key, deviceRaw := c.First(); key != nil; key, deviceRaw = c.Next() {
			var device devices.Device
			err := gob.NewDecoder(bytes.NewBuffer(deviceRaw)).Decode(&device)
			if err != nil {
				return errors.Wrap(err, "error problem to decoding the device struct")
			}

			devicesList = append(devicesList, device)
		}

		return nil
	})

	return devicesList, err
}

// GetXDevices returns X (set by count) devices starting at the firstDeviceUUID
func (ds DeviceStore) GetXDevices(firstDeviceUUID *string, count int64) ([]devices.Device, error) {
	var devicesList []devices.Device
	err := ds.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(devicesBucket)
		if bucket == nil {
			return errors.New("error in DeviceStore.GetXDevices: devices bucket does not exist")
		}

		c := bucket.Cursor()
		scanning := false
		scanCount := count
		for key, deviceRaw := c.First(); key != nil; key, deviceRaw = c.Next() {
			if !scanning && (firstDeviceUUID == nil || string(key) == *firstDeviceUUID) {
				scanning = true
			}

			if scanning {
				scanCount = scanCount - 1
				var device devices.Device
				err := gob.NewDecoder(bytes.NewBuffer(deviceRaw)).Decode(&device)
				if err != nil {
					return errors.Wrap(err, "error problem to decoding the device struct")
				}

				devicesList = append(devicesList, device)
			}

			if scanCount == 0 {
				break
			}
		}

		return nil
	})

	return devicesList, err

}

// Get returns a device from its UUID
func (ds DeviceStore) Get(uuid string) (devices.Device, error) {
	var device devices.Device
	err := ds.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(devicesBucket)
		if bucket == nil {
			return errors.New("error devices bucket does not exist")
		}

		deviceRaw := bucket.Get([]byte(uuid))
		if deviceRaw == nil {
			return errors.New("device not found")
		}

		err := gob.NewDecoder(bytes.NewBuffer(deviceRaw)).Decode(&device)

		return err
	})

	return device, err
}

// Search returns a list of device from a query
func (ds DeviceStore) Search(query string) ([]devices.Device, error) {
	var devicesList []devices.Device
	err := ds.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(devicesBucket)
		if bucket == nil {
			return errors.New("error in DeviceStore.Search: devices bucket does not exist")
		}

		c := bucket.Cursor()
		for key, deviceRaw := c.First(); key != nil; key, deviceRaw = c.Next() {
			var device devices.Device
			err := gob.NewDecoder(bytes.NewBuffer(deviceRaw)).Decode(&device)
			if err != nil {
				return errors.Wrap(err, "error problem to decoding the device struct")
			}

			if strings.Contains(device.UUID, query) || strings.Contains(device.DisplayName, query) || strings.Contains(device.Windows.DeviceID, query) || strings.Contains(device.Hardware.ID, query) {
				devicesList = append(devicesList, device)
			}
		}

		return nil
	})

	return devicesList, err
}

// EditOrCreate adds a new device or edits the existing device in the DB
func (ds DeviceStore) EditOrCreate(device devices.Device) error {
	// Encode Device
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(device); err != nil {
		return errors.Wrap(err, "error problem to encoding devices struct")
	}
	deviceRaw := buf.Bytes()

	// Store to DB
	err := ds.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(devicesBucket)
		if bucket == nil {
			return errors.New("error devices bucket does not exist")
		}

		err := bucket.Put([]byte(device.UUID), deviceRaw)
		return err
	})

	return err
}

// NewDeviceStore creates and initialises a new DeviceStore from a DB connection
func NewDeviceStore(db *bolt.DB) (DeviceStore, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(devicesBucket)
		return err
	})

	return DeviceStore{
		db,
	}, err
}
