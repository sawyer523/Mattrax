package enrollprovision

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/mdm/windows/protocol/wstep"
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
	"github.com/pkg/errors"
)

func Handler(server *mattrax.Server) http.HandlerFunc {
	managementServerURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.Domain,
		Path:   "/ManagementServer/MDM.svc",
	}).String()

	return func(w http.ResponseWriter, r *http.Request) {
		// Decode request from client
		var cmd Request
		if err := xml.NewDecoder(r.Body).Decode(&cmd); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify request structure
		if err := cmd.Verify(server.Config); err != nil {
			log.Println(errors.Wrap(err, "invalid MdeDiscoveryRequest:"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// TODO: Create Device in DB
		// defer func() {
		// 	// TODO: Save Device assuming no error occured
		// }()

		// Sign client CSR
		signedCertCommonName := "oscar@otbeaumont.me"         // TODO: Get User principal name
		if cmd.GetContextItem("EnrollmentType") == "Device" { // TODO: Possibly error if no EnrollmentType??
			signedCertCommonName = cmd.GetContextItem("DeviceID") // TODO: Possibly error if no DeviceID??
		}

		signedClientCert, clientCert, err := wstep.SignRequest(server.Certificates, cmd.Body.BinarySecurityToken.Value, signedCertCommonName)
		if err != nil {
			panic(err) // TODO
		}

		h := sha1.New()
		h.Write(signedClientCert)
		signedClientCertFingerprint := strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil))) // TODO: Cleanup

		// Prepare root identity cert details
		h2 := sha1.New()
		h2.Write(server.Certificates.IdentityCertRaw)
		identityCertFingerprint := strings.ToUpper(fmt.Sprintf("%x", h2.Sum(nil))) // TODO: Cleanup

		// Determain Certstore
		certStore := "User"
		if cmd.GetContextItem("EnrollmentType") == "Device" { // TODO: Possibly error no EnrollmentType??
			certStore = "System"
		}

		_ = clientCert // TEMP

		DMCLientProviderParameters := []WapParameter{}

		if server.Settings.TenantSupportPhone != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, WapParameter{
				Name:     "HelpPhoneNumber",
				Value:    server.Settings.TenantSupportPhone,
				DataType: "string",
			})
		}

		if server.Settings.TenantSupportWebsite != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, WapParameter{
				Name:     "HelpWebsite",
				Value:    server.Settings.TenantSupportWebsite,
				DataType: "string",
			})
		} else {
			DMCLientProviderParameters = append(DMCLientProviderParameters, WapParameter{
				Name:     "HelpWebsite",
				Value:    "https://mattrax.otbeaumont.me",
				DataType: "string",
			})
		}

		if server.Settings.TenantSupportEmail != "" {
			DMCLientProviderParameters = append(DMCLientProviderParameters, WapParameter{
				Name:     "HelpEmailAddress",
				Value:    server.Settings.TenantSupportEmail,
				DataType: "string",
			})
		}

		// Create provisioning profile
		resProvisioningProfile := WapProvisioningDoc{
			Version: "1.1",
			Characteristic: []WapCharacteristic{
				WapCharacteristic{
					// Spec: https://docs.microsoft.com/en-us/windows/client-management/mdm/certificatestore-csp
					Type: "CertificateStore",
					Characteristics: []WapCharacteristic{
						WapCharacteristic{
							Type: "Root",
							Characteristics: []WapCharacteristic{
								WapCharacteristic{
									Type: "System",
									Characteristics: []WapCharacteristic{
										WapCharacteristic{
											Type: identityCertFingerprint,
											Params: []WapParameter{
												WapParameter{
													Name:  "EncodedCertificate",
													Value: base64.StdEncoding.EncodeToString(server.Certificates.IdentityCertRaw),
												},
											},
										},
									},
								},
							},
						},
						WapCharacteristic{
							Type: "My",
							Characteristics: []WapCharacteristic{
								WapCharacteristic{
									Type: certStore,
									Characteristics: []WapCharacteristic{
										WapCharacteristic{
											Type: signedClientCertFingerprint,
											Params: []WapParameter{
												WapParameter{
													Name:  "EncodedCertificate",
													Value: base64.StdEncoding.EncodeToString(signedClientCert),
												},
											},
										},
										WapCharacteristic{ // TODO: ?
											Type:   "PrivateKeyContainer",
											Params: []WapParameter{},
										},
									},
								},
							},
						},
					},
				},
				WapCharacteristic{
					// Spec: https://docs.microsoft.com/en-us/windows/client-management/mdm/w7-application-csp
					Type: "APPLICATION",
					Params: []WapParameter{
						WapParameter{
							Name:  "APPID",
							Value: "w7",
						},
						WapParameter{
							Name:  "PROVIDER-ID",
							Value: "MattraxMDM",
							// DataType: "sting", //TODO: On all
						},
						WapParameter{
							Name:  "ADDR",
							Value: managementServerURL,
						},
						WapParameter{
							Name:  "NAME",
							Value: server.Settings.TenantName,
						},
						WapParameter{
							Name: "BACKCOMPATRETRYDISABLED",
						},
						WapParameter{
							Name:  "CONNRETRYFREQ",
							Value: "3",
						},
						WapParameter{
							Name:  "DEFAULTENCODING",
							Value: "application/vnd.syncml.dm+xml",
						},
						WapParameter{
							Name:  "INITIALBACKOFFTIME",
							Value: "16000", // Note: In milliseconds
						},
						WapParameter{
							Name:  "MAXBACKOFFTIME",
							Value: "86400000", // Note: In milliseconds
						},
						WapParameter{
							Name:  "PROTOVER",
							Value: "1.2",
						},
						WapParameter{
							Name:  "ROLE",
							Value: "4294967295", // TODO
						},
						// WapParameter{
						// 	Name:  "SSLCLIENTCERTSEARCHCRITERIA",
						// 	Value: "Subject=" + strings.ReplaceAll(url.PathEscape(clientCert.Subject.CommonName), "=", "%3D") + "&Stores=My%5C" + certStore,
						// 	// Value: "Subject=" + strings.ReplaceAll(url.PathEscape(clientCert.Subject.String()), "=", "%3D") + "&amp;Stores=My%5C" + certStore,
						// },
					},
					Characteristics: []WapCharacteristic{
						// TODO: APPAUTH
						WapCharacteristic{
							Type: "APPAUTH",
							Params: []WapParameter{
								WapParameter{
									Name:  "AAUTHLEVEL",
									Value: "CLIENT",
								},
								WapParameter{
									Name:  "AAUTHTYPE",
									Value: "DIGEST",
								},
								WapParameter{
									Name:  "AAUTHSECRET",
									Value: "dummy",
								},
								WapParameter{
									Name:  "AAUTHDATA",
									Value: "nonce",
								},
							},
						},
						WapCharacteristic{
							Type: "APPAUTH",
							Params: []WapParameter{
								WapParameter{
									Name:  "AAUTHLEVEL",
									Value: "APPSRV",
								},
								WapParameter{
									Name:  "AAUTHTYPE",
									Value: "DIGEST",
								},
								WapParameter{
									Name:  "AAUTHNAME",
									Value: "dummy",
								},
								WapParameter{
									Name:  "AAUTHSECRET",
									Value: "dummy",
								},
								WapParameter{
									Name:  "AAUTHDATA",
									Value: "nonce",
								},
							},
						},
					},
				},
				WapCharacteristic{
					Type: "DMClient",
					Characteristics: []WapCharacteristic{
						WapCharacteristic{
							Type: "Provider",
							Characteristics: []WapCharacteristic{
								WapCharacteristic{
									Type: "MattraxMDM",
									Params: append([]WapParameter{
										WapParameter{
											Name:     "EntDeviceName",
											Value:    "Demo Persons Device", // TODO: Device Name from Context
											DataType: "string",
										},
										WapParameter{
											Name:     "EntDMID",
											Value:    "aaaaaaa", // TODO: Mattrax DB Device ID
											DataType: "string",
										},
										// WapParameter{
										// 	Name:     "SignedEntDMID",
										// 	Value:    "", // TODO
										// 	DataType: "string",
										// },
										// WapParameter{
										// 	Name:     "CertRenewTimeStamp",
										// 	Value:    "", // TODO
										// 	DataType: "",
										// },
										// WapParameter{
										// 	Name:     "UPN",
										// 	Value:    "", // TODO: Email from user
										// 	DataType: "string",
										// },
										// WapParameter{
										// 	Name:     "RequireMessageSigning",
										// 	Value:    "true",
										// 	DataType: "boolean",
										// },
										// WapParameter{
										// 	Name:     "SyncApplicationVersion",
										// 	Value:    "2.0", // TODO: Is this correct 2.0
										// 	DataType: "string",
										// },
										// WapParameter{
										// 	Name:     "AADResourceID",
										// 	Value:    "", // TODO: Fill value
										// 	DataType: "string",
										// },
										WapParameter{
											Name:     "NumberOfDaysAfterLostContactToUnenroll",
											Value:    "730", // 2 years
											DataType: "integer",
										},
										// Note Mattrax currently doesn't support: ExchangeID, PublisherDeviceID, CommercialID
									}, DMCLientProviderParameters...),
									Characteristics: []WapCharacteristic{
										WapCharacteristic{
											Type: "Poll",
											Params: []WapParameter{
												WapParameter{
													Name:     "IntervalForFirstSetOfRetries",
													Value:    "15",
													DataType: "integer",
												},
												WapParameter{
													Name:     "NumberOfFirstRetries",
													Value:    "5",
													DataType: "integer",
												},
												WapParameter{
													Name:     "IntervalForSecondSetOfRetries",
													Value:    "60",
													DataType: "integer",
												},
												WapParameter{
													Name:     "NumberOfSecondRetries",
													Value:    "10",
													DataType: "integer",
												},
												WapParameter{
													Name:     "IntervalForRemainingScheduledRetries",
													Value:    "1440",
													DataType: "integer",
												},
												WapParameter{
													Name:     "NumberOfRemainingScheduledRetries",
													Value:    "0",
													DataType: "integer",
												},
												WapParameter{
													Name:     "PollOnLogin",
													Value:    "true",
													DataType: "boolean",
												},
												WapParameter{
													Name:     "AllUsersPollOnFirstLogin",
													Value:    "true",
													DataType: "boolean",
												},
											},
										},
										WapCharacteristic{
											Type: "CustomEnrollmentCompletePage",
											Params: []WapParameter{
												WapParameter{
													Name:     "Title",
													Value:    "Enrollment Complete",
													DataType: "string",
												},
												WapParameter{
													Name:     "BodyText",
													Value:    "Your device is now being managed by '" + server.Settings.TenantName + "'. Please contact your IT administrators for support if you have any problems.",
													DataType: "string",
												},
											},
										},
										// FUTURE: FirstSyncStatus like Apple DEP
										// WapCharacteristic{
										// 	Type: "EnhancedAppLayerSecurity",
										// 	Params: []WapParameter{
										// 		WapParameter{
										// 			Name:     "",
										// 			Value:    "",
										// 			DataType: "string",
										// 		},
										// 	},
										// },
									},
								},
							},
						},
					},
				},
			},
		}

		// Marshal provisioning profile
		provisioningProfile, err := xml.Marshal(resProvisioningProfile)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Create response
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
					Value:        base64.StdEncoding.EncodeToString(provisioningProfile),
				},
				RequestID: 0,
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
