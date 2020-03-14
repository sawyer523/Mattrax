package winmdm

import (
	"net/url"

	"github.com/mattrax/Mattrax/pkg/winmdm/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
)

// ProvisionRequest is the SOAP Envelope sent from the client during the Discovery phase
type ProvisionRequest struct {
	XMLName xml.Name             `xml:"s:Envelope"`
	Header  soap.Header          `xml:"s:Header"`
	Body    ProvisionRequestBody `xml:"s:Body>wst:RequestSecurityToken"`
}

// Verify checks the ProvisionRequest contains valid values and if not it returns a soap fault
func (cmd ProvisionRequest) Verify() (bool, soap.FaultEnvelop) {
	if cmd.Header.Action != "http://schemas.microsoft.com/windows/pki/2009/01/enrollment/RST/wstep" {
		return false, soap.NewBasicFault("s:Sender", "a:ActionMismatch", "client request body is invalid")
	}

	if url, err := url.ParseRequestURI(cmd.Header.To); err != nil || url.Scheme != "https" {
		return false, soap.NewEnrollmentFault("s:Sender", "a:EndpointUnavailable", "invalid To address", "EnrollmentServer", "EnrollmentInternalServiceError", "")
	}

	// TODO

	return true, soap.FaultEnvelop{}
}

// ProvisionRequestBody contains body of the request
type ProvisionRequestBody struct {
	TokenType           string                              `xml:"wst:TokenType"`
	RequestType         string                              `xml:"wst:RequestType"`
	BinarySecurityToken ProvisionRequestBinarySecurityToken `xml:"wsse:BinarySecurityToken"`
	AdditionalContext   []ProvisionRequestContextItem       `xml:"ac:AdditionalContext>ac:ContextItem"`
}

// ProvisionRequestBinarySecurityToken contains the CSR from the enrolling client
type ProvisionRequestBinarySecurityToken struct {
	ValueType    string `xml:"ValueType,attr"`
	EncodingType string `xml:"EncodingType,attr"`
	Value        string `xml:",chardata"`
}

// ProvisionRequestContextItem contains a key/value pair which contains information about the enrolling device
type ProvisionRequestContextItem struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"ac:Value"`
}

// GetAdditionalContextItem retrieves the first AdditionalContext item with the specified name
func (cmd RequestBody) GetAdditionalContextItem(name string) string {
	for _, contextItem := range cmd.AdditionalContext {
		if contextItem.Name == name {
			return contextItem.Value
		}
	}
	return ""
}

// GetAdditionalContextItems retrieves the AdditionalContext items with the specified name
func (cmd RequestBody) GetAdditionalContextItems(name string) []string {
	var contextItems []string
	for _, contextItem := range cmd.AdditionalContext {
		if contextItem.Name == name {
			contextItems = append(contextItems, contextItem.Value)
		}
	}
	return contextItems
}
