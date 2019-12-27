package enrollprovision

import (
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
)

type HeaderSecurityTimestamp struct {
	ID      string `xml:"u:Id,attr"`
	Created string `xml:"u:Created"`
	Expires string `xml:"u:Expires"`
}

type HeaderSecurity struct {
	NamespaceO     string                  `xml:"xmlns:o,attr"`
	MustUnderstand string                  `xml:"s:mustUnderstand,attr"`
	Timestamp      HeaderSecurityTimestamp `xml:"u:Timestamp"`
}

type BinarySecurityToken struct {
	ValueType    string `xml:"ValueType,attr"`
	EncodingType string `xml:"EncodingType,attr"`
	Value        string `xml:",chardata"`
}

type ResponseBody struct {
	TokenType           string              `xml:"RequestSecurityTokenResponse>TokenType"`
	DispositionMessage  string              `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollment RequestSecurityTokenResponse>DispositionMessage"` // TODO: Invalid type
	BinarySecurityToken BinarySecurityToken `xml:"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd RequestSecurityTokenResponse>RequestedSecurityToken>BinarySecurityToken"`
	RequestID           int                 `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollment RequestSecurityTokenResponse>RequestID"`
}

type ResponseEnvelope struct {
	XMLName      xml.Name            `xml:"s:Envelope"`
	NamespaceS   string              `xml:"xmlns:s,attr"`
	NamespaceA   string              `xml:"xmlns:a,attr"`
	NamespaceU   string              `xml:"xmlns:u,attr"`
	HeaderAction soap.MustUnderstand `xml:"s:Header>a:Action"`
	// HeaderActivityID string                   `xml:"s:Header>a:ActivityID"` // TODO: Is this needed
	HeaderRelatesTo string         `xml:"s:Header>a:RelatesTo"`
	HeaderSecurity  HeaderSecurity `xml:"s:Header>o:Security"`
	Body            ResponseBody   `xml:"http://docs.oasis-open.org/ws-sx/ws-trust/200512 s:Body>RequestSecurityTokenResponseCollection"`
}
