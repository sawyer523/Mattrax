package enrolldiscovery

import (
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
)

type Response struct {
	AuthPolicy                 string `xml:"DiscoverResult>AuthPolicy"`
	EnrollmentVersion          string `xml:"DiscoverResult>EnrollmentVersion"`
	EnrollmentPolicyServiceURL string `xml:"DiscoverResult>EnrollmentPolicyServiceUrl"`
	EnrollmentServiceURL       string `xml:"DiscoverResult>EnrollmentServiceUrl"`
	AuthenticationServiceUrl   string `xml:"DiscoverResult>AuthenticationServiceUrl,omitempty"`
}

type ResponseBody struct {
	NamespaceXSI     string   `xml:"xmlns:xsi,attr"`
	NamespaceXSD     string   `xml:"xmlns:xsd,attr"`
	DiscoverResponse Response `xml:"http://schemas.microsoft.com/windows/management/2012/01/enrollment DiscoverResponse"`
}

type ResponseEnvelope struct {
	XMLName          xml.Name            `xml:"s:Envelope"`
	NamespaceS       string              `xml:"xmlns:s,attr"`
	NamespaceA       string              `xml:"xmlns:a,attr"`
	HeaderAction     soap.MustUnderstand `xml:"s:Header>a:Action"`
	HeaderActivityID string              `xml:"s:Header>a:ActivityID"`
	HeaderRelatesTo  string              `xml:"s:Header>a:RelatesTo"`
	Body             ResponseBody        `xml:"s:Body"`
}
