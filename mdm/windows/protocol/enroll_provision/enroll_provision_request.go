package enrollprovision

import (
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
)

// Request contains the SOAP request Envelope
type Request struct {
	XMLName xml.Name    `xml:"s:Envelope"`
	Header  soap.Header `xml:"s:Header"`
	Body    RequestBody `xml:"s:Body>wst:RequestSecurityToken"`
}

// RequestBody contains the body of the SOAP Envelope
type RequestBody struct {
	TokenType           string                     `xml:"wst:TokenType"`
	RequestType         string                     `xml:"wst:RequestType"`
	BinarySecurityToken RequestBinarySecurityToken `xml:"wsse:BinarySecurityToken"`
	AdditionalContext   []ContextItem              `xml:"ac:AdditionalContext>ac:ContextItem"`
}

// RequestBinarySecurityToken contains the CSR from the enrolling client
type RequestBinarySecurityToken struct {
	ValueType    string `xml:"ValueType,attr"`
	EncodingType string `xml:"EncodingType,attr"`
	Value        string `xml:",chardata"`
}

// ContextItem contains a key/value pair which contains information about the enrolling device
type ContextItem struct {
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
