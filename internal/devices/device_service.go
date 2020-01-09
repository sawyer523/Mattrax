package devices

// Service contains the code for interfacing with devices.
type Service interface {
	GetAll() ([]Device, error)
	GetXDevices(firstDeviceUUID *string, count int64) ([]Device, error)
	Get(uuid string) (Device, error)
	Search(query string) ([]Device, error)
	EditOrCreate(device Device) error
}
