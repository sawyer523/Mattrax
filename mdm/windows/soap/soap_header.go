package soap

// Header is the SOAP Header for a request. It contains the intent, id and authentication details for a SOAP request.
type Header struct {
	Action    string `xml:"a:Action"`
	MessageID string `xml:"a:MessageID,omitempty"`
	ReplyTo   struct {
		Address string `xml:"a:Address"`
	} `xml:"a:ReplyTo,omitempty"`
	To           string                `xml:"a:To,omitempty"`
	WSSESecurity HeaderMdeWSSESecurity `xml:"wsse:Security,omitempty"`
}

// HeaderRes is the SOAP Header for a response. It contains the intent, request id and devices request id for a SOAP response.
type HeaderRes struct {
	Action     MustUnderstand `xml:"a:Action"`
	ActivityID string         `xml:"a:ActivityID"`
	RelatesTo  string         `xml:"a:RelatesTo"`
}

// HeaderMdeWSSESecurity is contained in SOAP Header and it carries the user authentication details.
type HeaderMdeWSSESecurity struct {
	Username            string `xml:"wsse:UsernameToken>wsse:Username"`
	Password            string `xml:"wsse:UsernameToken>wsse:Password"`
	BinarySecurityToken string `xml:"wsse:BinarySecurityToken"`
}
