package windows

import (
	"net/http"
	"strconv"

	"github.com/mattrax/Mattrax/pkg/winmdm"
	"github.com/mattrax/Mattrax/pkg/winmdm/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
	"github.com/rs/zerolog/log"
)

func (mdm *MDM) DiscoveryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			return
		}

		var cmd winmdm.DiscoveryRequest
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)
		if err := xml.NewDecoder(r.Body).Decode(&cmd); err != nil {
			log.Debug().Str("type", "error").Str("remote-addr", r.RemoteAddr).Err(err).Msg("error: discovery request: failed to parse client request body")
			fault := soap.NewBasicFault("s:Sender", "s:MessageFormat", "client request body couldn't be parsed")
			fault.Response(w)
			return
		}
		if ok, fault := cmd.Verify("Federated"); !ok {
			fault.Response(w)
			return
		}

		res := cmd.NewResponse(winmdm.DiscoveryResponseBody{
			AuthPolicy:                 "Federated",
			EnrollmentPolicyServiceURL: mdm.EnrollmentPolicyServiceURL,
			EnrollmentServiceURL:       mdm.EnrollmentProvisionServiceURL,
			AuthenticationServiceURL:   mdm.AuthenticationServiceURL,
		})

		// Marshal and send the response to client
		response, err := xml.Marshal(res)
		if err != nil {
			fault := soap.NewBasicFault("s:Reciever", "a:InternalServiceFault", "Mattrax error: failed to generate discovery response")
			fault.Response(w)
			return
		}
		w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
		_, err = w.Write(response)
		if err != nil {
			log.Error().Str("email", cmd.Body.EmailAddress).Str("device-type", cmd.Body.DeviceType).Err(err).Msg("error: failed to send discovery response body")
		}
	}
}
