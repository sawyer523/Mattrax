package enrolldiscovery

import (
	"net/http"
	"net/url"
	"strconv"
	"sync"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/generic"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
	"github.com/rs/zerolog/log"
)

// GETHandler handles the HTTP GET request for discovery.
// The handler returns a HTTP status 200 so the device can determine if an enrollment server exists
// It MUST be mounted at the path "/EnrollmentServer/Discovery.svc"
func GETHandler(server *mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

// Handler handles the HTTP POST request for discovery.
// The handler parses the device details and responds with the location of the enrollment endpoints and the authentication policy required
// It MUST be mounted at the path "/EnrollmentServer/Discovery.svc"
func Handler(server *mattrax.Server) http.HandlerFunc {
	const maxRequestBodySize = 5000

	enrollmentPolicyServiceURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.Domain,
		Path:   "/EnrollmentServer/Policy.svc",
	}).String()

	enrollmentServiceURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.Domain,
		Path:   "/EnrollmentServer/Enrollment.svc",
	}).String()

	federationServiceURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.Domain,
		Path:   "/EnrollmentServer/Authenticate",
	}).String()

	var (
		authPolicy = "Federated" // For now this is the only supported by Mattrax
		once       sync.Once
		response   []byte
	)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > maxRequestBodySize {
			fault := soap.NewBasicFault("s:Sender", "s:MessageFormat", "client request body too large to process")
			fault.Response(w)
			return
		}

		var cmd Request
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
		if err := xml.NewDecoder(r.Body).Decode(&cmd); err != nil {
			log.Debug().Str("type", "error").Str("remote-addr", r.RemoteAddr).Err(err).Msg("error: discovery request: failed to parse client request body")
			fault := soap.NewBasicFault("s:Sender", "s:MessageFormat", "client request body couldn't be parsed")
			fault.Response(w)
			return
		}

		if cmd.Header.Action != "http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/Discover" {
			fault := soap.NewBasicFault("s:Sender", "a:ActionMismatch", "client request body is invalid")
			fault.Response(w)
			return
		}

		if url, err := url.ParseRequestURI(cmd.Header.To); err != nil || url.Scheme != "https" {
			fault := soap.NewEnrollmentFault("s:Sender", "a:EndpointUnavailable", "invalid To address", "EnrollmentServer", "EnrollmentInternalServiceError", "")
			fault.Response(w)
			return
		}

		if cmd.Body.EmailAddress == "" || !types.ValidEmail.MatchString(cmd.Body.EmailAddress) {
			fault := soap.NewBasicFault("s:Sender", "s:MessageFormat", "client request body missing or invalid email address")
			fault.Response(w)
			return
		}

		// Note: Intune doesn't verify this but I though it would prevent issues after enrollment if not supported.
		if cmd.Body.RequestVersion != "4.0" {
			fault := soap.NewEnrollmentFault("s:Sender", "a:InternalServiceFault", "the discovery request version "+cmd.Body.RequestVersion+" is not supported", "DeviceNotSupported", "unsupported discovery request version", "")
			fault.Response(w)
			return
		}

		if cmd.Body.DeviceType != "CIMClient_Windows" {
			fault := soap.NewEnrollmentFault("s:Sender", "a:InternalServiceFault", "device type "+cmd.Body.DeviceType+" is not supported", "DeviceNotSupported", "unsupported device type", "")
			fault.Response(w)
			return
		}

		// Note: Intune disregards what the device supports and returns the AuthPolicy it desires
		if !cmd.Body.AuthPolicies.IsAuthPolicySupported(authPolicy) {
			fault := soap.NewEnrollmentFault("s:Sender", "a:InternalServiceFault", authPolicy+" auth policy is not supported by your device by required for enrollment", "DeviceNotSupported", "unsupported auth policy", "")
			fault.Response(w)
			return
		}

		once.Do(func() {
			res := Response{
				NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
				NamespaceA: "http://www.w3.org/2005/08/addressing",
				Header: soap.HeaderRes{
					Action: soap.MustUnderstand{
						MustUnderstand: "1",
						Value:          "http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/DiscoverResponse",
					},
					ActivityID: generic.GenerateID(),
					RelatesTo:  cmd.Header.MessageID,
				},
				Body: ResponseBody{
					NamespaceXSI: "http://www.w3.org/2001/XMLSchema-instance",
					NamespaceXSD: "http://www.w3.org/2001/XMLSchema",
					DiscoverResponse: DiscoverResponse{
						AuthPolicy:                 authPolicy,
						EnrollmentVersion:          cmd.Body.RequestVersion,
						EnrollmentPolicyServiceURL: enrollmentPolicyServiceURL,
						EnrollmentServiceURL:       enrollmentServiceURL,
						AuthenticationServiceURL:   federationServiceURL,
					},
				},
			}

			// Marshal and send the response to client
			var err error
			response, err = xml.Marshal(res)
			if err != nil {
				fault := soap.NewBasicFault("s:Reciever", "a:InternalServiceFault", "mattrax error: failed to generate discovery response")
				fault.Response(w)
				return
			}
		})

		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		_, err := w.Write(response)
		if err != nil {
			log.Error().Str("email", cmd.Body.EmailAddress).Str("device-type", cmd.Body.DeviceType).Err(err).Msg("error: failed to send discovery response body")
		}
	}
}
