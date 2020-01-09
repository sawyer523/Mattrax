package enrollprovision

import (
	"encoding/base64"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/certificates"
	"github.com/mattrax/Mattrax/internal/devices"
)

func GenerateProvisioningProfile(server *mattrax.Server, managementServerURL string, identityCertificate certificates.Identity, device devices.Device, clientCertificateDer []byte) WapProvisioningDoc {
	certStore := "User"
	if device.Windows.EnrollmentType == "Device" { // TODO: Possibly error no EnrollmentType??
		certStore = "System"
	}

	serverSettings := server.Settings.Get()

	DMCLientProviderParameters := []WapParameter{}

	if serverSettings.Tenant.SupportPhone != "" {
		DMCLientProviderParameters = append(DMCLientProviderParameters, WapParameter{
			Name:     "HelpPhoneNumber",
			Value:    serverSettings.Tenant.SupportPhone,
			DataType: "string",
		})
	}

	if serverSettings.Tenant.SupportEmail != "" {
		DMCLientProviderParameters = append(DMCLientProviderParameters, WapParameter{
			Name:     "HelpEmailAddress",
			Value:    serverSettings.Tenant.SupportEmail,
			DataType: "string",
		})
	}

	if serverSettings.Tenant.SupportWebsite != "" {
		DMCLientProviderParameters = append(DMCLientProviderParameters, WapParameter{
			Name:     "HelpWebsite",
			Value:    serverSettings.Tenant.SupportWebsite,
			DataType: "string",
		})
	} else {
		DMCLientProviderParameters = append(DMCLientProviderParameters, WapParameter{
			Name:     "HelpWebsite",
			Value:    "https://mattrax.otbeaumont.me",
			DataType: "string",
		})
	}

	return WapProvisioningDoc{
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
										Type: identityCertificate.CertHash,
										Params: []WapParameter{
											WapParameter{
												Name:  "EncodedCertificate",
												Value: base64.StdEncoding.EncodeToString(identityCertificate.CertRaw), // TODO: Can CertRaw be removed by using .Cert.Raw ???
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
										Type: device.IdentityCertificate.Hash,
										Params: []WapParameter{
											WapParameter{
												Name:  "EncodedCertificate",
												Value: base64.StdEncoding.EncodeToString(clientCertificateDer),
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
						Value: serverSettings.Tenant.Name,
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
												Value:    "Your device is now being managed by '" + serverSettings.Tenant.Name + "'. Please contact your IT administrators for support if you have any problems.",
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
}
