package protocol

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	wtypes "github.com/mattrax/Mattrax/mdm/windows/types"
	"github.com/mattrax/Mattrax/pkg/xml"
	perrors "github.com/pkg/errors"
)

func Discover(server mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func Discovery(server mattrax.Server) http.HandlerFunc {
	enrollmentPolicyServiceURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.PrimaryDomain,
		Path:   "/EnrollmentServer/Policy.svc",
	}).String()

	enrollmentServiceURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.PrimaryDomain,
		Path:   "/EnrollmentServer/Enrollment.svc",
	}).String()

	return func(w http.ResponseWriter, r *http.Request) {
		// Verify client user-agent
		if r.Header.Get("User-Agent") != "ENROLLClient" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Decode request from client
		var cmd wtypes.MdeDiscoveryRequest
		if err := xml.NewDecoder(r.Body).Decode(&cmd); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify request structure
		if err := cmd.VerifyStructure(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get Matrax settings
		settings, err := server.SettingsService.Get()
		if err != nil {
			log.Println(perrors.Wrap(err, "error DiscoveryPOST: failed to retrieve settings"))

		}

		// Verify request
		if err := cmd.VerifyContext(server.Config, settings); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create response
		var authPolicy string
		var authenticationServiceUrl string
		if settings.Windows.AuthPolicy == types.AuthPolicyOnPremise {
			if !cmd.IsAuthPolicySupport(types.AuthPolicyOnPremise) {
				log.Println(errors.New("error DiscoveryPOST: device doesn't support OnPremise AuthPolicy but is required by server"))
				w.WriteHeader(http.StatusConflict)
				return
			}
			authPolicy = "OnPremise"
			authenticationServiceUrl = ""
		} else if settings.Windows.AuthPolicy == types.AuthPolicyFederated {
			if !cmd.IsAuthPolicySupport(types.AuthPolicyFederated) {
				log.Println(errors.New("error DiscoveryPOST: device doesn't support Federated AuthPolicy but is required by server"))
				w.WriteHeader(http.StatusConflict)
				return
			}
			authPolicy = "Federated"
			authenticationServiceUrl = settings.Windows.FederationPortalURL
		} else {
			log.Println(errors.New("error DiscoveryPOST: invalid AuthPolicy '" + string(settings.Windows.AuthPolicy) + "'"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := wtypes.MdeDiscoveryResponseEnvelope{
			NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
			NamespaceA: "http://www.w3.org/2005/08/addressing",
			HeaderAction: wtypes.MustUnderstand{
				MustUnderstand: "1",
				Value:          "http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/DiscoverResponse",
			},
			HeaderActivityID: wtypes.GenerateActivityID(),
			HeaderRelatesTo:  cmd.Header.MessageID,
			Body: wtypes.MdeDiscoveryResponseBody{
				NamespaceXSI: "http://www.w3.org/2001/XMLSchema-instance",
				NamespaceXSD: "http://www.w3.org/2001/XMLSchema",
				DiscoverResponse: wtypes.DiscoverResponse{
					AuthPolicy:                 authPolicy,
					EnrollmentVersion:          cmd.Body.Discover.Request.RequestVersion,
					EnrollmentPolicyServiceURL: enrollmentPolicyServiceURL,
					EnrollmentServiceURL:       enrollmentServiceURL,
					AuthenticationServiceUrl:   authenticationServiceUrl,
				},
			},
		}

		// Marshal and send the response to client
		if response, err := xml.Marshal(res); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			w.Write(response)
		}
	}
}
