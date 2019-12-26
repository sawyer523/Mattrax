package wtypes

import (
	"errors"
	"net/url"
	"strings"

	"github.com/mattrax/Mattrax/pkg/xml"
)

type MdeEnrollmentRequest struct {
	XMLName xml.Name `xml:"s:Envelope"`
	Header  struct {
		Action    string `xml:"a:Action"`
		MessageID string `xml:"a:MessageID"`
		ReplyTo   struct {
			Address string `xml:"a:Address"`
		} `xml:"a:ReplyTo"`
		To           string          `xml:"a:To"`
		WSSESecurity MdeWSSESecurity `xml:"wsse:Security"`
	} `xml:"s:Header"`
	Body struct {
		TokenType           string `xml:"wst:RequestSecurityToken>wst:TokenType"`
		RequestType         string `xml:"wst:RequestSecurityToken>wst:RequestType"`
		BinarySecurityToken struct {
			ValueType    string `xml:"ValueType,attr"`
			EncodingType string `xml:"EncodingType,attr"`
			Value        string `xml:",chardata"`
		} `xml:"wst:RequestSecurityToken>wsse:BinarySecurityToken"`
		AdditionalContext struct {
			ContextItems []struct {
				Name  string `xml:"Name,attr"`
				Value string `xml:"ac:Value"`
			} `xml:"ac:ContextItem"`
		} `xml:"wst:RequestSecurityToken>ac:AdditionalContext"`
	} `xml:"s:Body"`
}

// VerifyStructure checks that the cmd contains a valid structure.
// This is done before more resource intensive checks to reduce the DoS vector.
func (cmd MdeEnrollmentRequest) VerifyStructure() error {
	if cmd.Header.Action == "" {
		return errors.New("invalid MdeEnrollmentRequest: empty Action")
	} else if cmd.Header.Action != "http://schemas.microsoft.com/windows/pki/2009/01/enrollment/RST/wstep" {
		return errors.New("invalid MdeEnrollmentRequest: invalid Action '" + cmd.Header.Action + "'")
	}

	if cmd.Header.MessageID == "" {
		return errors.New("invalid MdeEnrollmentRequest: empty MessageID")
	} else if !strings.HasPrefix(cmd.Header.MessageID, "urn:uuid:") {
		return errors.New("invalid MdeEnrollmentRequest: invalid MessageID prefix '" + cmd.Header.MessageID + "'")
	} else if !validMessageID.MatchString(cmd.Header.MessageID) {
		return errors.New("invalid MdeEnrollmentRequest: invalid characters in MessageID '" + cmd.Header.MessageID + "'")
	}

	if cmd.Header.To == "" {
		return errors.New("invalid MdeEnrollmentRequest: empty To")
	} else if _, err := url.ParseRequestURI(cmd.Header.To); err != nil {
		return errors.New("invalid MdeEnrollmentRequest: invalid To '" + cmd.Header.To + "'")
	}

	if cmd.Header.WSSESecurity.BinarySecurityToken != "" { // Federated Authentication
		if !binarySecurityToken.MatchString(cmd.Header.WSSESecurity.BinarySecurityToken) {
			return errors.New("invalid MdeEnrollmentRequest: invalid binary security token")
		}
	} else if cmd.Header.WSSESecurity.Username != "" { // On Premise Authentication
		if !validEmail.MatchString(cmd.Header.WSSESecurity.Username) {
			// Note: the incorrect username is not displayed as it could have been a user accidentally typing thier password
			return errors.New("invalid MdeEnrollmentRequest: invalid email address")
		}

		if cmd.Header.WSSESecurity.Password == "" {
			return errors.New("invalid MdeEnrollmentRequest: empty password")
		} else if !validPassword.MatchString(cmd.Header.WSSESecurity.Password) {
			// Note: the incorrect password is not displayed for what should be obvious reasons
			return errors.New("invalid MdeEnrollmentRequest: invalid password")
		}
	} else {
		return errors.New("invalid MdeEnrollmentRequest: no supported client authentication details received")
	}

	if len(cmd.Body.BinarySecurityToken.Value) == 0 {
		return errors.New("invalid MdeEnrollmentRequest: invalid body binary security token")
	}

	return nil
}

type MdeWapProvisioningDoc struct {
	XMLName        xml.Name               `xml:"wap-provisioningdoc"`
	Version        string                 `xml:"version,attr"`
	Characteristic []MdeWapCharacteristic `xml:"characteristic"`
}

type MdeWapCharacteristic struct {
	Type           string `xml:"type,attr,omitempty"`
	Params         []MdeWapParm
	Characteristic []MdeWapCharacteristic `xml:"characteristic,omitempty"` // TODO: rename to Characteristics
}

type MdeWapParm struct {
	XMLName  xml.Name `xml:"parm"`
	Name     string   `xml:"name,attr,omitempty"`
	Value    string   `xml:"value,attr,omitempty"`
	DataType string   `xml:"datatype,attr,omitempty"`
}

type MdeEnrollmentHeaderSecurityTimestamp struct {
	ID      string `xml:"u:Id,attr"`
	Created string `xml:"u:Created"`
	Expires string `xml:"u:Expires"`
}

type MdeEnrollmentHeaderSecurity struct {
	NamespaceO     string                               `xml:"xmlns:o,attr"`
	MustUnderstand string                               `xml:"s:mustUnderstand,attr"`
	Timestamp      MdeEnrollmentHeaderSecurityTimestamp `xml:"u:Timestamp"`
}

type MdeBinarySecurityToken struct {
	ValueType    string `xml:"ValueType,attr"`
	EncodingType string `xml:"EncodingType,attr"`
	Value        string `xml:",chardata"`
}

type MdeEnrollmentResponseBody struct {
	TokenType           string                 `xml:"RequestSecurityTokenResponse>TokenType"`
	DispositionMessage  string                 `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollment RequestSecurityTokenResponse>DispositionMessage"` // TODO: Invalid type
	BinarySecurityToken MdeBinarySecurityToken `xml:"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd RequestSecurityTokenResponse>RequestedSecurityToken>BinarySecurityToken"`
	RequestID           int                    `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollment RequestSecurityTokenResponse>RequestID"`
}

// MdeEnrollmentResponseEnvelope contains the payload sent from the server to the client telling it // TODO: Finish this from spec
type MdeEnrollmentResponseEnvelope struct {
	XMLName      xml.Name       `xml:"s:Envelope"`
	NamespaceS   string         `xml:"xmlns:s,attr"`
	NamespaceA   string         `xml:"xmlns:a,attr"`
	NamespaceU   string         `xml:"xmlns:u,attr"`
	HeaderAction MustUnderstand `xml:"s:Header>a:Action"`
	// HeaderActivityID string                   `xml:"s:Header>a:ActivityID"` // TODO: Is this needed
	HeaderRelatesTo string                      `xml:"s:Header>a:RelatesTo"`
	HeaderSecurity  MdeEnrollmentHeaderSecurity `xml:"s:Header>o:Security"`
	Body            MdeEnrollmentResponseBody   `xml:"http://docs.oasis-open.org/ws-sx/ws-trust/200512 s:Body>RequestSecurityTokenResponseCollection"`
}
