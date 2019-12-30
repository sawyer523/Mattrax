package types

// MDMProtcol defines the protocol a device is being managed by
type MDMProtcol int

// DeviceUUID is a unique identifier given to each device by the MDM server upon enrollment
type DeviceUUID []byte

const (
	// Windows is MDMProtocol used to manage Windows devices
	Windows MDMProtcol = iota
	// Apple is the MDMProtocol used to manage Apple devices
	Apple
)

// Device is an electronic device that is managed by the MDM server
type Device struct {
	UUID            DeviceUUID
	DisplayName     string
	Owner           User
	SerialNumber    string // This is set by the devices manufacturer
	DeviceModel     string // This is set by the devices manufacturer
	OperatingSystem string // This is set by the devices manufacturer
	IMEI            string // May be blank
	MEID            string // May be blank
	MDMProtcol      MDMProtcol
	Policies        []Policy
}

// DeviceService contains the implemented functionality for devices
type DeviceService interface {
	GetAll() ([]Device, error)
	Get(UUID DeviceUUID) (Device, error)
	EnrollOrEdit(UUID DeviceUUID, device Device) error
	GenerateUUID() (DeviceUUID, error)
}
