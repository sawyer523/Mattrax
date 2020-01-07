package enrolldiscovery

import (
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
)

// Response contains the SOAP response Envelope
type Response struct {
	XMLName          xml.Name            `xml:"s:Envelope"`
	NamespaceS       string              `xml:"xmlns:s,attr"`
	NamespaceA       string              `xml:"xmlns:a,attr"`
	HeaderAction     soap.MustUnderstand `xml:"s:Header>a:Action"`
	HeaderActivityID string              `xml:"s:Header>a:ActivityID"`
	HeaderRelatesTo  string              `xml:"s:Header>a:RelatesTo"`
	Body             ResponseBody        `xml:"s:Body"`
}

// ResponseBody contains the body of the SOAP Envelope
type ResponseBody struct {
	NamespaceXSI     string           `xml:"xmlns:xsi,attr"`
	NamespaceXSD     string           `xml:"xmlns:xsd,attr"`
	DiscoverResponse DiscoverResponse `xml:"http://schemas.microsoft.com/windows/management/2012/01/enrollment DiscoverResponse"`
}

// DiscoverResponse contains the enrollment endpoints and authentication policy for the device to continue enrollment with
type DiscoverResponse struct {
	AuthPolicy                 string `xml:"DiscoverResult>AuthPolicy"`
	EnrollmentVersion          string `xml:"DiscoverResult>EnrollmentVersion"`
	EnrollmentPolicyServiceURL string `xml:"DiscoverResult>EnrollmentPolicyServiceUrl"`
	EnrollmentServiceURL       string `xml:"DiscoverResult>EnrollmentServiceUrl"`
	AuthenticationServiceURL   string `xml:"DiscoverResult>AuthenticationServiceUrl,omitempty"`
}
