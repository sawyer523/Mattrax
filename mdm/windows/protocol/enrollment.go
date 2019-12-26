package protocol

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
	wtypes "github.com/mattrax/Mattrax/mdm/windows/types"
	"github.com/mattrax/Mattrax/pkg/xml"
)

func Enrollment(server mattrax.Server) http.HandlerFunc {
	managementServerURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.PrimaryDomain,
		Path:   "/ManagementServer/MDM.svc",
	}).String()

	managementServerListURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.PrimaryDomain,
		Path:   "/ManagementServer/ServerList.svc",
	}).String()

	return func(w http.ResponseWriter, r *http.Request) {
		// Verify client user-agent
		if r.Header.Get("User-Agent") != "ENROLLClient" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Decode request from client
		var cmd wtypes.MdeEnrollmentRequest
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

		// TODO: VerifyContext

		// TODO: Create Device in DB
		// defer func() {
		// 	// TODO: Save Device assuming no error occured
		// }()

		// Sign client CSR
		signedClientCert, clientCert, err := server.CertificateService.SignWSTEPRequest(cmd.Body.BinarySecurityToken.Value)
		if err != nil {
			panic(err) // TODO
		}

		h := sha1.New()
		h.Write(signedClientCert)
		signedClientCertFingerprint := strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil))) // TODO: Cleanup

		// Prepare root identity cert details
		identityCertificateRaw, _, err := server.CertificateService.GetIdentityRaw()
		if err != nil {
			panic(err) // TODO
		}

		h2 := sha1.New()
		h2.Write(identityCertificateRaw)
		identityCertFingerprint := strings.ToUpper(fmt.Sprintf("%x", h2.Sum(nil))) // TODO: Cleanup

		// Get MDM settings
		settings, err := server.SettingsService.Get()
		if err != nil {
			panic(err) // TODO
		}

		// Create provisioning profile
		resProvisioningProfile := wtypes.MdeWapProvisioningDoc{
			Version: "1.1",
			Characteristic: []wtypes.MdeWapCharacteristic{
				wtypes.MdeWapCharacteristic{
					Type: "CertificateStore",
					Characteristic: []wtypes.MdeWapCharacteristic{
						wtypes.MdeWapCharacteristic{
							Type: "Root",
							Characteristic: []wtypes.MdeWapCharacteristic{
								wtypes.MdeWapCharacteristic{
									Type: "System",
									Characteristic: []wtypes.MdeWapCharacteristic{
										wtypes.MdeWapCharacteristic{
											Type: identityCertFingerprint,
											Params: []wtypes.MdeWapParm{
												wtypes.MdeWapParm{
													Name:  "EncodedCertificate",
													Value: base64.StdEncoding.EncodeToString(identityCertificateRaw),
												},
											},
										},
									},
								},
							},
						},
						wtypes.MdeWapCharacteristic{
							Type: "My",
							Characteristic: []wtypes.MdeWapCharacteristic{
								wtypes.MdeWapCharacteristic{
									Type: "User",
									Characteristic: []wtypes.MdeWapCharacteristic{
										wtypes.MdeWapCharacteristic{
											Type: signedClientCertFingerprint,
											Params: []wtypes.MdeWapParm{
												wtypes.MdeWapParm{
													Name:  "EncodedCertificate",
													Value: base64.StdEncoding.EncodeToString(signedClientCert),
												},
											},
										},
										wtypes.MdeWapCharacteristic{
											Type: "PrivateKeyContainer",
											Params: []wtypes.MdeWapParm{
												wtypes.MdeWapParm{
													Name:  "KeySpec",
													Value: "2",
												},
												wtypes.MdeWapParm{
													Name:  "ContainerName",
													Value: "ConfigMgrEnrollment",
												},
												wtypes.MdeWapParm{
													Name:  "ProviderType",
													Value: "1",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				wtypes.MdeWapCharacteristic{
					Type: "APPLICATION",
					Params: []wtypes.MdeWapParm{
						wtypes.MdeWapParm{
							Name:  "APPID",
							Value: "w7",
						},
						wtypes.MdeWapParm{
							Name:  "PROVIDER-ID",
							Value: "Mattrax MDM Server",
						},
						wtypes.MdeWapParm{
							Name:  "NAME",
							Value: settings.TenantName,
						},
						wtypes.MdeWapParm{
							Name:  "SSPHyperlink",
							Value: "http://go.microsoft.com/fwlink/?LinkId=255310", // Enterprise Management App
						},
						wtypes.MdeWapParm{
							Name:  "ADDR",
							Value: managementServerURL,
						},
						wtypes.MdeWapParm{
							Name:  "ServerList",
							Value: managementServerListURL,
						},
						wtypes.MdeWapParm{
							Name:  "ROLE",
							Value: "4294967295", // ? Possible Values
						},
						/* Discriminator to set whether the client should do Certificate Revocation List checking. */
						wtypes.MdeWapParm{
							Name:  "CRLCheck",
							Value: "0",
						},
						wtypes.MdeWapParm{
							Name:  "CONNRETRYFREQ",
							Value: "6",
						},
						wtypes.MdeWapParm{
							Name:  "INITIALBACKOFFTIME",
							Value: "30000",
						},
						wtypes.MdeWapParm{
							Name:  "MAXBACKOFFTIME",
							Value: "120000",
						},
						wtypes.MdeWapParm{
							Name: "BACKCOMPATRETRYDISABLED",
						},
						wtypes.MdeWapParm{
							Name:  "DEFAULTENCODING",
							Value: "application/vnd.syncml.dm+wbxml",
						},
						// TODO: This is causing issues
						// wtypes.MdeWapParm{
						// 	Name:  "SSLCLIENTCERTSEARCHCRITERIA",
						// 	Value: "Subject=CN=%3d" + clientCert.Subject.CommonName + "&amp;Stores=MY%5CUser", // TODO: Correct Value
						// },
					},
					Characteristic: []wtypes.MdeWapCharacteristic{
						wtypes.MdeWapCharacteristic{
							Type: "APPAUTH",
							Params: []wtypes.MdeWapParm{
								wtypes.MdeWapParm{
									Name:  "AAUTHLEVEL",
									Value: "CLIENT",
								},
								wtypes.MdeWapParm{
									Name:  "AAUTHTYPE",
									Value: "DIGEST",
								},
								wtypes.MdeWapParm{
									Name:  "AAUTHSECRET",
									Value: "dummy",
								},
								wtypes.MdeWapParm{
									Name:  "AAUTHDATA",
									Value: "nonce",
								},
							},
						},
						wtypes.MdeWapCharacteristic{
							Type: "APPAUTH",
							Params: []wtypes.MdeWapParm{
								wtypes.MdeWapParm{
									Name:  "AAUTHLEVEL",
									Value: "APPSRV",
								},
								wtypes.MdeWapParm{
									Name:  "AAUTHTYPE",
									Value: "DIGEST",
								},
								wtypes.MdeWapParm{
									Name:  "AAUTHNAME",
									Value: "dummy",
								},
								wtypes.MdeWapParm{
									Name:  "AAUTHSECRET",
									Value: "dummy",
								},
								wtypes.MdeWapParm{
									Name:  "AAUTHDATA",
									Value: "nonce",
								},
							},
						},
					},
				},
				wtypes.MdeWapCharacteristic{
					Type: "Registry",
					Characteristic: []wtypes.MdeWapCharacteristic{
						wtypes.MdeWapCharacteristic{
							Type: `HKLM\Security\MachineEnrollment`,
							Params: []wtypes.MdeWapParm{
								wtypes.MdeWapParm{
									Name:     "RenewalPeriod",
									Value:    "363",
									DataType: "integer",
								},
							},
						},
						wtypes.MdeWapCharacteristic{
							Type: `HKLM\Security\MachineEnrollment\OmaDmRetry`,
							Params: []wtypes.MdeWapParm{
								wtypes.MdeWapParm{
									Name:     "NumRetries",
									Value:    "8",
									DataType: "integer",
								},
								wtypes.MdeWapParm{
									Name:     "RetryInterval",
									Value:    "15",
									DataType: "integer",
								},
								wtypes.MdeWapParm{
									Name:     "AuxNumRetries",
									Value:    "5",
									DataType: "integer",
								},
								wtypes.MdeWapParm{
									Name:     "AuxRetryInterval",
									Value:    "3",
									DataType: "integer",
								},
								wtypes.MdeWapParm{
									Name:     "Aux2NumRetries",
									Value:    "0",
									DataType: "integer",
								},
								wtypes.MdeWapParm{
									Name:     "Aux2RetryInterval",
									Value:    "480",
									DataType: "integer",
								},
							},
						},
					},
				},
				wtypes.MdeWapCharacteristic{
					Type: "Registry",
					Characteristic: []wtypes.MdeWapCharacteristic{
						wtypes.MdeWapCharacteristic{
							Type: `HKLM\Software\Windows\CurrentVersion\MDM\MachineEnrollment`,
							Params: []wtypes.MdeWapParm{
								wtypes.MdeWapParm{
									Name:     "DeviceName",
									Value:    "TODO",
									DataType: "string",
								},
							},
						},
					},
				},
				wtypes.MdeWapCharacteristic{
					Type: "Registry",
					Characteristic: []wtypes.MdeWapCharacteristic{
						wtypes.MdeWapCharacteristic{
							Type: `HKLM\SOFTWARE\Windows\CurrentVersion\MDM\MachineEnrollment`,
							Params: []wtypes.MdeWapParm{
								// Thumbprint of root certificate.
								wtypes.MdeWapParm{
									Name:     "SslServerRootCertHash",
									Value:    identityCertFingerprint,
									DataType: "string",
								},
								// Store for device certificate.
								wtypes.MdeWapParm{
									Name:     "SslClientCertStore",
									Value:    "MY%5CSystem",
									DataType: "string",
								},
								wtypes.MdeWapParm{
									Name:     "SslClientCertSubjectName",
									Value:    "Subject=CN=%3d" + clientCert.Subject.CommonName, //"CN%3de4c6b893-07a7-4b24-878e-9d8602c3d289",
									DataType: "string",
								},
								wtypes.MdeWapParm{
									Name:     "SslClientCertHash",
									Value:    signedClientCertFingerprint,
									DataType: "string",
								},
							},
						},
						wtypes.MdeWapCharacteristic{
							Type: `HKLM\Security\Provisioning\OMADM\Accounts\037B1F0D3842015588E753CDE76EC724`,
							Params: []wtypes.MdeWapParm{
								wtypes.MdeWapParm{
									Name:     "SslClientCertReference",
									Value:    "My;System;" + signedClientCertFingerprint,
									DataType: "string",
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

		fmt.Println(string(provisioningProfile)) // TEMP

		// Create response
		res := wtypes.MdeEnrollmentResponseEnvelope{
			NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
			NamespaceA: "http://www.w3.org/2005/08/addressing",
			NamespaceU: "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd",
			HeaderAction: wtypes.MustUnderstand{
				MustUnderstand: "1",
				Value:          "http://schemas.microsoft.com/windows/pki/2009/01/enrollment/RSTRC/wstep",
			},
			HeaderRelatesTo: cmd.Header.MessageID,
			HeaderSecurity: wtypes.MdeEnrollmentHeaderSecurity{
				NamespaceO:     "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd",
				MustUnderstand: "1",
				Timestamp: wtypes.MdeEnrollmentHeaderSecurityTimestamp{
					// TODO: all these values do what??
					ID:      "_0",
					Created: "2018-11-30T00:32:59.420Z",
					Expires: "2018-12-30T00:37:59.420Z",
				},
			},
			Body: wtypes.MdeEnrollmentResponseBody{
				TokenType:          "http://schemas.microsoft.com/5.0.0.0/ConfigurationManager/Enrollment/DeviceEnrollmentToken",
				DispositionMessage: "", // TODO: Wrong type + What does it do?
				BinarySecurityToken: wtypes.MdeBinarySecurityToken{
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
