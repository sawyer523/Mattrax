package protocol

import (
	"log"
	"net/http"
	"strconv"

	mattrax "github.com/mattrax/Mattrax/internal"
	wtypes "github.com/mattrax/Mattrax/mdm/windows/types"
	"github.com/mattrax/Mattrax/pkg/xml"
)

func Policy(server mattrax.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify client user-agent
		if r.Header.Get("User-Agent") != "ENROLLClient" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Decode request from client
		var cmd wtypes.MdeGetPoliciesRequest
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

		// Verify request
		if err := cmd.VerifyContext(server.Config, server.UserService); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			// TODO: Return error that user can understand
			return
		}

		// TODO: Use the Request Body to work out if changes have happened for PolciiesNotChanged. What is requestFilter???
		// TODO: This response is hardcoded. Automatically generate from CertService + Settings

		res := wtypes.MdePolicyResponseEnvelope{
			NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
			NamespaceA: "http://www.w3.org/2005/08/addressing",
			NamespaceU: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
			HeaderAction: wtypes.MustUnderstand{
				MustUnderstand: "1",
				Value:          "http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy/IPolicy/GetPoliciesResponse",
			},
			// HeaderActivityID: wtypes.GenerateActivityID(),
			HeaderRelatesTo: cmd.Header.MessageID,
			Body: wtypes.MdePolicyResponseBody{
				NamespaceXSI: "http://www.w3.org/2001/XMLSchema-instance",
				NamespaceXSD: "http://www.w3.org/2001/XMLSchema",
				PoliciesResponse: wtypes.MdePoliciesResponse{
					PolicyID:           wtypes.GenerateID(), // TODO: Does this have to stay constant????
					PolicyFriendlyName: "Mattrax Identity",  // TODO: Does it show
					NextUpdateHours:    12,                  // TODO: After 12 hours does it request this same endpoint like Apple?????????
					PoliciesNotChanged: false,               // TODO: Track this. False means policies have changed since last updateHour
					Policies: []wtypes.MdePolicy{
						wtypes.MdePolicy{
							OIDReference: 0,
							CAs: wtypes.MdeCACollection{
								Nil: true,
							},
							Attributes: wtypes.MdeAttributes{
								CommonName:   "Mattrax Identity2", // TODO: Does it show
								PolicySchema: 3,
								CertificateValidity: wtypes.MdeAttributesCertificateValidity{
									// TODO: What is good for these values. Also how does renewal work??
									ValidityPeriodSeconds: 1209600,
									RenewalPeriodSeconds:  172800,
								},
								EnrollmentPermission: wtypes.MdeEnrollmentPermission{
									Enroll:     true,  // TODO: Try false as rejection
									AutoEnroll: false, // TODO: See what is changes
								},
								PrivateKeyAttributes: wtypes.MdePrivateKeyAttributes{
									MinimalKeyLength: 2048, // TODO: Get from CertService
									KeySpec: wtypes.MdeKeySpec{
										Nil: true,
									},
									KeyUsageProperty: wtypes.MdeKeyUsageProperty{
										Nil: true,
									},
									Permissions: wtypes.MdePermissions{
										Nil: true,
									},
									AlgorithmOIDReference: wtypes.MdeAlgorithmOIDReference{
										Nil: true,
									},
									CryptoProviders: wtypes.MdeCryptoProviders{
										Nil: true,
									},
								},
								Revision: wtypes.MdeRevision{
									MajorRevision: 101, // TODO: Change and see what happens. Version control inside Mattrax???
									MinorRevision: 0,   // TODO: Change and see what happens.
								},
								SupersededPolicies: wtypes.MdeSupersededPolicies{
									Nil: true,
								},
								PrivateKeyFlags: wtypes.MdePrivateKeyFlags{
									Nil: true,
								},
								SubjectNameFlags: wtypes.MdeSubjectNameFlags{
									Nil: true,
								},
								EnrollmentFlags: wtypes.MdeEnrollmentFlags{
									Nil: true,
								},
								GeneralFlags: wtypes.MdeGeneralFlags{
									Nil: true,
								},
								HashAlgorithmOIDReference: 0, // TODO: What dis do?
								RARequirements: wtypes.MdeRARequirements{
									Nil: true,
								},
								KeyArchivalAttributes: wtypes.MdeKeyArchivalAttributes{
									Nil: true,
								},
								Extensions: wtypes.MdeExtensions{
									Nil: true,
								},
							},
						},
					},
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
