package devices

import (
	"errors"

	"github.com/google/go-cmp/cmp"
	"github.com/imdario/mergo"
	"github.com/rs/zerolog/log"
	"github.com/samsarahq/thunder/graphql/schemabuilder"
)

// MountAPI attaches the Devices Schema to the GraphQL API
func MountAPI(s Service, builder *schemabuilder.Schema) {
	deviceObject := builder.Object("Device", Device{})
	deviceObject.Description = "A device is an electronic device that is managed by the MDM server"

	var mdmProtocolEnum MDMProtcol
	builder.Enum(mdmProtocolEnum, map[string]MDMProtcol{
		"Windows": WindowsMDM,
		"Apple":   AppleMDM,
	})

	identityObject := builder.Object("DeviceIdentityCertificate", DeviceIdentityCertificate{})
	// identityObject.FieldFunc("subject", func() string { return "TODO" })

	// TOOD: Does this work
	identityObject.FieldFunc("subject", func(identityCertificate DeviceIdentityCertificate) string {
		return identityCertificate.Subject.String()
	})

	query := builder.Query()
	query.FieldFunc("devices", func(req struct {
		FirstDevice *string
		Count       int64
	}) ([]Device, error) {
		if req.FirstDevice == nil && req.Count == int64(0) {
			return s.GetAll()
		}

		return s.GetXDevices(req.FirstDevice, req.Count)
	})
	query.FieldFunc("getDevice", func(req struct {
		UUID string `graphql:"uuid"`
	}) (Device, error) {
		if req.UUID != "" {
			return s.Get(req.UUID)
		}

		return Device{}, errors.New("invalid request: no device identifier was given")
	})
	query.FieldFunc("searchDevices", func(req struct{ Query string }) ([]Device, error) {
		if req.Query == "" {
			return nil, errors.New("invalid request: no query was given")
		}
		return s.Search(req.Query)
	})

	mutation := builder.Mutation()
	mutation.FieldFunc("updateDevice", func(newDevice Device) (Device, error) {
		currentDevice, err := s.Get(newDevice.UUID)
		if err != nil {
			return Device{}, err
		}
		if err := mergo.Merge(&newDevice, currentDevice); err != nil {
			log.Error().Err(err).Msg("error merging Settings structs")
			return Device{}, errors.New("internal server error: failed to merge settings")
		}

		// Device UUID is read only through the API.
		if newDevice.UUID != currentDevice.UUID {
			return Device{}, errors.New("the device uuid is read only")
		}

		// Enrollment details are read only through the API.
		if newDevice.EnrolledAt != currentDevice.EnrolledAt || !cmp.Equal(newDevice.EnrolledBy, currentDevice.EnrolledBy) {
			return Device{}, errors.New("the device's enrollment details are read only")
		}

		if err := s.EditOrCreate(newDevice); err != nil {
			return Device{}, err
		}
		return newDevice, nil
	})
}
