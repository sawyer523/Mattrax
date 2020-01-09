package devices

import (
	"crypto/x509/pkix"
	"time"

	"github.com/mattrax/Mattrax/internal/types"
)

// MDMProtcol defines the protocol a device is being managed by
type MDMProtcol int

const (
	// WindowsMDM is MDMProtocol used to manage Windows devices
	WindowsMDM MDMProtcol = iota
	// AppleMDM is the MDMProtocol used to manage Apple devices
	AppleMDM
)

// Device is an electronic device that is managed by the MDM server
// TODO: Make cross platform. Currently it has lots of Windows only values.
type Device struct {
	UUID                string                    `graphql:"uuid"`      // A unique identifier given to each device by the MDM server upon enrollment
	DisplayName         string                    `graphql:",optional"` // The user friendly name for the device.
	Protocol            MDMProtcol                `graphql:",optional"` // The MDM Protocol that manages the device
	EnrolledAt          time.Time                 `graphql:",optional"` // Time device was enrolled in MDM (Read only)
	EnrolledBy          types.User                `graphql:",optional"` // The user that enrolled the device in MDM (Stores UUID only in struct as reference) (Read only)
	Windows             WindowsDevice             `graphql:",optional"`
	Hardware            DeviceHardware            `graphql:",optional"`
	IdentityCertificate DeviceIdentityCertificate `graphql:",optional"`

	// TODO: Policies        []Policy
}

// DeviceHardware contains details about the physical device managed by MDM
type DeviceHardware struct {
	ID  string   `graphql:",optional"` // HardwareID
	MAC []string `graphql:",optional"`
}

// DeviceIdentityCertificate contains detials about the identity certificate issued to the device
type DeviceIdentityCertificate struct {
	Subject   pkix.Name `graphql:"-"`
	Hash      string    `graphql:",optional"`
	NotBefore time.Time `graphql:",optional"`
	NotAfter  time.Time `graphql:",optional"`
}

// TODO: move to Windows package
type WindowsDevice struct {
	DeviceID           string `graphql:",optional"`
	DeviceType         string `graphql:",optional"`
	EnrollmentType     string `graphql:",optional"`
	OSEdition          string `graphql:",optional"`
	OSVersion          string `graphql:",optional"`
	ApplicationVersion string `graphql:",optional"`
}
