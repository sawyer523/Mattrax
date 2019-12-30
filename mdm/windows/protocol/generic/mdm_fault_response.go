package generic

import "github.com/mattrax/Mattrax/pkg/xml"

type SoapFaultSubcodeValue struct {
	NamespaceA string `xml:"xmlns:a,attr"`
	Value      string `xml:",innerxml"`
}

type SoapFaultReasonText struct {
	Language string `xml:"xml:lang,attr"`
	Value    string `xml:",innerxml"`
}

type SoapFaultEnvelop struct {
	XMLName      xml.Name              `xml:"s:Envelope"`
	NamespaceS   string                `xml:"xmlns:s,attr"`
	NamespaceA   string                `xml:"xmlns:a,attr"`
	Value        string                `xml:"s:Body>s:Fault>s:Code>s:Value"`
	SubcodeValue SoapFaultSubcodeValue `xml:"s:Body>s:Fault>s:Code>s:Subcode>s:Value"`
	ReasonText   SoapFaultReasonText   `xml:"s:Body>s:Fault>s:Reason>s:Text"`
}

func NewGenericSoapFault(subcode string, msg string) SoapFaultEnvelop {
	return SoapFaultEnvelop{
		NamespaceS: "http://www.w3.org/2003/05/soap-envelope",
		NamespaceA: "http://www.w3.org/2005/08/addressing",
		Value:      "s:Receiver",
		SubcodeValue: SoapFaultSubcodeValue{
			NamespaceA: "http://schemas.microsoft.com/net/2005/12/windowscommunicationfoundation/dispatcher",
			Value:      subcode,
		},
		ReasonText: SoapFaultReasonText{
			Language: "en-US",
			Value:    "Mattrax error: " + msg,
		},
	}
}
