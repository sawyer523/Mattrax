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

	managementServerListURL := (&url.URL{
		Scheme: "https",
		Host:   server.Config.Domain,
		Path:   "/ManagementServer/ServerList.svc",
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
		signedClientCert, clientCert, err := wstep.SignRequest(server.Certificates, cmd.Body.BinarySecurityToken.Value)
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

		// Create provisioning profile
		resProvisioningProfile := WapProvisioningDoc{
			Version: "1.1",
			Characteristic: []WapCharacteristic{
				WapCharacteristic{
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
									Type: "User",
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
										WapCharacteristic{
											Type: "PrivateKeyContainer",
											Params: []WapParameter{
												WapParameter{
													Name:  "KeySpec",
													Value: "2",
												},
												WapParameter{
													Name:  "ContainerName",
													Value: "ConfigMgrEnrollment",
												},
												WapParameter{
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
				WapCharacteristic{
					Type: "APPLICATION",
					Params: []WapParameter{
						WapParameter{
							Name:  "APPID",
							Value: "w7",
						},
						WapParameter{
							Name:  "PROVIDER-ID",
							Value: "Mattrax MDM Server",
						},
						WapParameter{
							Name:  "NAME",
							Value: server.Settings.TenantName,
						},
						WapParameter{
							Name:  "SSPHyperlink",
							Value: "http://go.microsoft.com/fwlink/?LinkId=255310", // Enterprise Management App
						},
						WapParameter{
							Name:  "ADDR",
							Value: managementServerURL,
						},
						WapParameter{
							Name:  "ServerList",
							Value: managementServerListURL,
						},
						WapParameter{
							Name:  "ROLE",
							Value: "4294967295", // ? Possible Values
						},
						/* Discriminator to set whether the client should do Certificate Revocation List checking. */
						WapParameter{
							Name:  "CRLCheck",
							Value: "0",
						},
						WapParameter{
							Name:  "CONNRETRYFREQ",
							Value: "6",
						},
						WapParameter{
							Name:  "INITIALBACKOFFTIME",
							Value: "30000",
						},
						WapParameter{
							Name:  "MAXBACKOFFTIME",
							Value: "120000",
						},
						WapParameter{
							Name: "BACKCOMPATRETRYDISABLED",
						},
						WapParameter{
							Name:  "DEFAULTENCODING",
							Value: "application/vnd.syncml.dm+wbxml",
						},
						// TODO: This is causing issues
						// WapParameter{
						// 	Name:  "SSLCLIENTCERTSEARCHCRITERIA",
						// 	Value: "Subject=CN=%3d" + clientCert.Subject.CommonName + "&amp;Stores=MY%5CUser", // TODO: Correct Value
						// },
					},
					Characteristics: []WapCharacteristic{
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
					Type: "Registry",
					Characteristics: []WapCharacteristic{
						WapCharacteristic{
							Type: `HKLM\Security\MachineEnrollment`,
							Params: []WapParameter{
								WapParameter{
									Name:     "RenewalPeriod",
									Value:    "363",
									DataType: "integer",
								},
							},
						},
						WapCharacteristic{
							Type: `HKLM\Security\MachineEnrollment\OmaDmRetry`,
							Params: []WapParameter{
								WapParameter{
									Name:     "NumRetries",
									Value:    "8",
									DataType: "integer",
								},
								WapParameter{
									Name:     "RetryInterval",
									Value:    "15",
									DataType: "integer",
								},
								WapParameter{
									Name:     "AuxNumRetries",
									Value:    "5",
									DataType: "integer",
								},
								WapParameter{
									Name:     "AuxRetryInterval",
									Value:    "3",
									DataType: "integer",
								},
								WapParameter{
									Name:     "Aux2NumRetries",
									Value:    "0",
									DataType: "integer",
								},
								WapParameter{
									Name:     "Aux2RetryInterval",
									Value:    "480",
									DataType: "integer",
								},
							},
						},
					},
				},
				WapCharacteristic{
					Type: "Registry",
					Characteristics: []WapCharacteristic{
						WapCharacteristic{
							Type: `HKLM\Software\Windows\CurrentVersion\MDM\MachineEnrollment`,
							Params: []WapParameter{
								WapParameter{
									Name:     "DeviceName",
									Value:    "TODO",
									DataType: "string",
								},
							},
						},
					},
				},
				WapCharacteristic{
					Type: "Registry",
					Characteristics: []WapCharacteristic{
						WapCharacteristic{
							Type: `HKLM\SOFTWARE\Windows\CurrentVersion\MDM\MachineEnrollment`,
							Params: []WapParameter{
								// Thumbprint of root certificate.
								WapParameter{
									Name:     "SslServerRootCertHash",
									Value:    identityCertFingerprint,
									DataType: "string",
								},
								// Store for device certificate.
								WapParameter{
									Name:     "SslClientCertStore",
									Value:    "MY%5CSystem",
									DataType: "string",
								},
								WapParameter{
									Name:     "SslClientCertSubjectName",
									Value:    "Subject=CN=%3d" + clientCert.Subject.CommonName, //"CN%3de4c6b893-07a7-4b24-878e-9d8602c3d289",
									DataType: "string",
								},
								WapParameter{
									Name:     "SslClientCertHash",
									Value:    signedClientCertFingerprint,
									DataType: "string",
								},
							},
						},
						WapCharacteristic{
							Type: `HKLM\Security\Provisioning\OMADM\Accounts\037B1F0D3842015588E753CDE76EC724`,
							Params: []WapParameter{
								WapParameter{
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
			/* TEMP */
			// response = []byte(`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing"><s:Header><a:Action s:mustunderstand="1">http://schemas.microsoft.com/windows/pki/2009/01/enrollment/rstrc/wstep</a:Action><ActivityID xmlns="http://schemas.microsoft.com/2004/09/servicemodel/diagnostics">` + generic.GenerateID() + `</ActivityID><a:RelatesTo>` + cmd.Header.Action + `</a:RelatesTo></s:Header><s:Body><s:Fault><s:Code><s:Value>s:receiver</s:value><s:Subcode><s:Value>s:Authorization</s:Value></s:Subcode></s:Code><s:Reason><s:Text xml:lang="en-US">This User is not authorized to enroll</s:text></s:Reason></s:Fault></s:Body></s:Envelope>`)
			// response = []byte(`<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing">
			// 	<s:Header>
			// 		<a:Action s:mustunderstand="1">http://schemas.microsoft.com/windows/pki/2009/01/enrollment/IWindowsDeviceEnrollmentService/RequestSecurityTokenWindowsDeviceEnrollmentServiceErrorFault</a:Action>
			// 		<ActivityID xmlns="http://schemas.microsoft.com/2004/09/servicemodel/diagnostics">` + generic.GenerateID() + `</ActivityID>
			// 		<a:RelatesTo>` + cmd.Header.Action + `</a:RelatesTo>
			// 	</s:Header>
			// 	<s:Body>
			// 	<s:fault>
			// 	<s:code>
			// 		<s:value>s:receiver</s:value>
			// 		<s:subcode>
			// 			<s:value>s:authorization</s:value>
			// 		</s:subcode>
			// 	</s:code>
			// 	<s:reason>
			// 		<s:text xml:lang="en-us">device cap reached</s:text>
			// 	</s:reason>
			// 	<s:detail>
			// 		<deviceenrollmentserviceerror xmlns="http://schemas.microsoft.com/windows/pki/2009/01/enrollment">
			// 			<errortype>devicecapreached</errortype>
			// 			<message>device cap reached</message>
			// 			<traceid>2493ee37-beeb-4cb9-833c-cadde9067645</traceid>
			// 		</deviceenrollmentserviceerror>
			// 	</s:detail>
			// </s:fault>

			// 	</s:Body>
			// </s:Envelope>`)
			/* END TEMP */

			w.Header().Set("Content-Type", "application/soap+xml; charset=utf-8")
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			w.Write(response)
		}
	}
}
