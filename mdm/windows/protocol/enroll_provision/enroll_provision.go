package enrollprovision

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/devices"
	"github.com/mattrax/Mattrax/internal/generic"
	"github.com/mattrax/Mattrax/internal/settings"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
	"github.com/rs/zerolog/log"
)

func Handler(server *mattrax.Server) http.HandlerFunc {
	const maxRequestBodySize = 5000

	managementServerURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.Domain,
		Path:   "/ManagementServer/Manage.svc",
	}).String()

	return func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > maxRequestBodySize {
			fault := soap.NewBasicFault("s:Sender", "s:MessageFormat", "client request body too large to process")
			fault.Response(w)
			return
		}

		if server.Settings.Get().ServerState != settings.StateNormal {
			// TODO: Propper Advanced Fault
			fault := soap.NewBasicFault("s:Sender", "a:EndpointUnavailable", "the server is currently not ready to accept enrollments")
			fault.Response(w)
			return
		}

		var cmd Request
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
		if err := xml.NewDecoder(r.Body).Decode(&cmd); err != nil {
			log.Debug().Str("type", "error").Str("remote-addr", r.RemoteAddr).Err(err).Msg("error: provision request: failed to parse client request body")
			fault := soap.NewBasicFault("s:Sender", "s:MessageFormat", "client request body couldn't be parsed")
			fault.Response(w)
			return
		}

		if cmd.Header.Action != "http://schemas.microsoft.com/windows/pki/2009/01/enrollment/RST/wstep" {
			fault := soap.NewBasicFault("s:Sender", "a:ActionMismatch", "client request body is invalid")
			fault.Response(w)
			return
		}

		if url, err := url.ParseRequestURI(cmd.Header.To); err != nil || url.Scheme != "https" {
			fault := soap.NewEnrollmentFault("s:Sender", "a:EndpointUnavailable", "invalid To address", "EnrollmentServer", "EnrollmentInternalServiceError", "")
			fault.Response(w)
			return
		}

		fmt.Println(cmd.Header.WSSESecurity.BinarySecurityToken) // TODO: Verify (check valid Microsoft or Mattrax token)

		// FINISH CHECKING INPUT: Verify CSR exists
		// Required EnrollmentType and possibly DeviceID

		// fault := soap.NewEnrollmentFault("s:Receiver", "s:Authorization", "device cap reached", "DeviceCapReached", "device cap reached", "2493ee37-beeb-4cb9-833c-cadde9067645")
		// fault.Response(w)

		// fault := soap.NewEnrollmentFault("s:Receiver", "a:InternalServiceFault", "hello world", "NotSupported", "Device Not Supported", "test")
		// fault.Response(w)

		// Binary Security Token
		// Check if device exists. Get users other devices.
		// AAD isManaged true

		device := devices.Device{
			UUID:        generic.GenerateID(),
			DisplayName: cmd.Body.GetAdditionalContextItem("DeviceName"),
			Protocol:    devices.WindowsMDM,
			EnrolledAt:  time.Now(),
			EnrolledBy:  types.User{
				// TODO: UUID
			},
			Hardware: devices.DeviceHardware{
				ID:  cmd.Body.GetAdditionalContextItem("HWDevID"),
				MAC: cmd.Body.GetAdditionalContextItems("MAC"),
			},
			Windows: devices.WindowsDevice{
				DeviceID:           cmd.Body.GetAdditionalContextItem("DeviceID"),
				DeviceType:         cmd.Body.GetAdditionalContextItem("DeviceType"),
				EnrollmentType:     cmd.Body.GetAdditionalContextItem("EnrollmentType"),
				OSEdition:          cmd.Body.GetAdditionalContextItem("OSEdition"),
				OSVersion:          cmd.Body.GetAdditionalContextItem("OSVersion"),
				ApplicationVersion: cmd.Body.GetAdditionalContextItem("ApplicationVersion"),
			},
		}

		faultLogger := log.With().Str("device-uuid", device.UUID).Str("device-display-name", device.DisplayName).Str("win-device-id", device.Windows.DeviceID).Logger()

		// defer func() {
		// 	// server.DeviceService
		// }()

		identityCertificate := server.Certificates.Get().Identity

		clientCertificateDer, err := SignClientCertificate(server, device, cmd.Body.BinarySecurityToken.Value, identityCertificate)
		if err != nil {
			// TODO
			return
		}

		fmt.Println(device)

		provisioningProfile := GenerateProvisioningProfile(server, managementServerURL, identityCertificate, device, clientCertificateDer)

		provisioningProfileXML, err := xml.Marshal(provisioningProfile)
		if err != nil {
			log.Debug().Str("type", "error").Str("remote-addr", r.RemoteAddr).Err(err).Msg("error: discovery request: failed to parse client request body")
			fault := soap.NewBasicFault("s:Reciever", "a:InternalServiceFault", "mattrax error: failed to generate provisioning profile")
			fault.Response(w)
			return
		}

		res := ResponseEnvelope{
			NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
			NamespaceA: "http://www.w3.org/2005/08/addressing",
			NamespaceU: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
			HeaderAction: soap.MustUnderstand{
				MustUnderstand: "1",
				Value:          "http://schemas.microsoft.com/windows/pki/2009/01/enrollment/RSTRC/wstep",
			},
			HeaderRelatesTo: cmd.Header.MessageID,
			HeaderSecurity: HeaderSecurity{
				NamespaceO:     "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd",
				MustUnderstand: "1",
				Timestamp: HeaderSecurityTimestamp{
					// TODO: all these values do what??
					ID:      "_0",
					Created: "2018-11-30T00:32:59.420Z",
					Expires: "2018-12-30T00:37:59.420Z",
				},
			},
			Body: ResponseBody{
				TokenType:          "http://schemas.microsoft.com/5.0.0.0/ConfigurationManager/Enrollment/DeviceEnrollmentToken",
				DispositionMessage: "", // TODO: Wrong type + What does it do?
				BinarySecurityToken: BinarySecurityToken{
					ValueType:    "http://schemas.microsoft.com/5.0.0.0/ConfigurationManager/Enrollment/DeviceEnrollmentProvisionDoc",
					EncodingType: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd#base64binary",
					Value:        base64.StdEncoding.EncodeToString(provisioningProfileXML),
				},
				RequestID: 0,
			},
		}

		// Marshal and send the response to client
		response, err := xml.Marshal(res)
		if err != nil {
			fault := soap.NewBasicFault("s:Reciever", "a:InternalServiceFault", "mattrax error: failed to generate provision response")
			fault.Response(w)
			return
		}

		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		_, err = w.Write(response)
		if err != nil {
			faultLogger.Error().Err(err).Msg("error: failed to send provision response body")
		}
	}
}
