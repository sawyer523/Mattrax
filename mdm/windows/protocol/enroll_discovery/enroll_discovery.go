package enrolldiscovery

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	generic "github.com/mattrax/Mattrax/mdm/windows/protocol/generic"
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
	"github.com/pkg/errors"
)

func GETHandler(server mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func Handler(server mattrax.Server) http.HandlerFunc {
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

	internalFederationServiceURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.PrimaryDomain,
		Path:   "/EnrollmentServer/Authenticate",
	}).String()

	return func(w http.ResponseWriter, r *http.Request) {
		// Decode request from client
		var cmd Request
		if err := xml.NewDecoder(r.Body).Decode(&cmd); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		/* TEMP: Getting settings will be refactored to not be expensive so this is temporary */
		settings, err := server.SettingsService.Get()
		if err != nil {
			panic(err)
		}
		/* END TEMP */

		if err := cmd.Verify(server.Config, settings); err != nil {
			log.Println(errors.Wrap(err, "invalid MdeDiscoveryRequest:"))
			w.WriteHeader(http.StatusBadRequest)
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
			if settings.Windows.FederationPortalURL == "" {
				authenticationServiceUrl = internalFederationServiceURL

			} else {
				authenticationServiceUrl = settings.Windows.FederationPortalURL
			}
		} else {
			log.Println(errors.New("error DiscoveryPOST: invalid AuthPolicy '" + string(settings.Windows.AuthPolicy) + "'"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := ResponseEnvelope{
			NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
			NamespaceA: "http://www.w3.org/2005/08/addressing",
			HeaderAction: soap.MustUnderstand{
				MustUnderstand: "1",
				Value:          "http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/DiscoverResponse",
			},
			HeaderActivityID: generic.GenerateID(),
			HeaderRelatesTo:  cmd.Header.MessageID,
			Body: ResponseBody{
				NamespaceXSI: "http://www.w3.org/2001/XMLSchema-instance",
				NamespaceXSD: "http://www.w3.org/2001/XMLSchema",
				DiscoverResponse: Response{
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
