package winmdm

import (
	"github.com/mattrax/Mattrax/pkg/winmdm/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
	uuid "github.com/satori/go.uuid"
)

// DiscoveryResponse is the SOAP Envelope returned to the client during the Discovery phase
type DiscoveryResponse struct {
	XMLName    xml.Name                       `xml:"s:Envelope"`
	NamespaceS string                         `xml:"xmlns:s,attr"`
	NamespaceA string                         `xml:"xmlns:a,attr"`
	Header     soap.HeaderRes                 `xml:"s:Header"`
	Body       DiscoveryResponseBodyContainer `xml:"s:Body"`
}

// DiscoveryResponseBody contains body of the response
type DiscoveryResponseBodyContainer struct {
	NamespaceXSI     string                `xml:"xmlns:xsi,attr"`
	NamespaceXSD     string                `xml:"xmlns:xsd,attr"`
	DiscoverResponse DiscoveryResponseBody `xml:"http://schemas.microsoft.com/windows/management/2012/01/enrollment DiscoverResponse"`
}

// DiscoveryResponseBodyDiscoverResponse contains the servers endpoints so the device can begin its enrollment
type DiscoveryResponseBody struct {
	AuthPolicy                 string `xml:"DiscoverResult>AuthPolicy"`
	EnrollmentVersion          string `xml:"DiscoverResult>EnrollmentVersion"`
	EnrollmentPolicyServiceURL string `xml:"DiscoverResult>EnrollmentPolicyServiceUrl"`
	EnrollmentServiceURL       string `xml:"DiscoverResult>EnrollmentServiceUrl"`
	AuthenticationServiceURL   string `xml:"DiscoverResult>AuthenticationServiceUrl,omitempty"`
}

func (cmd DiscoveryRequest) NewResponse(body DiscoveryResponseBody) DiscoveryResponse {
	body.EnrollmentVersion = cmd.Body.RequestVersion
	return DiscoveryResponse{
		NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
		NamespaceA: "http://www.w3.org/2005/08/addressing",
		Header: soap.HeaderRes{
			Action: soap.MustUnderstand{
				MustUnderstand: "1",
				Value:          "http://schemas.microsoft.com/windows/management/2012/01/enrollment/IDiscoveryService/DiscoverResponse",
			},
			ActivityID: uuid.NewV4().String(),
			RelatesTo:  cmd.Header.MessageID,
		},
		Body: DiscoveryResponseBodyContainer{
			NamespaceXSI:     "http://www.w3.org/2001/XMLSchema-instance",
			NamespaceXSD:     "http://www.w3.org/2001/XMLSchema",
			DiscoverResponse: body,
		},
	}
}
