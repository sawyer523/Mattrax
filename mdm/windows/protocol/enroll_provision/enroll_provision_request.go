package enrollprovision

import (
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/mdm/windows/protocol/generic"
	"github.com/mattrax/Mattrax/pkg/xml"
	"github.com/pkg/errors"
)

type Request struct {
	XMLName xml.Name       `xml:"s:Envelope"`
	Header  generic.Header `xml:"s:Header"`
	Body    struct {
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

func (cmd Request) Verify(config mattrax.Config) error {
	/* Verify Structure: Lightwieght structure and datatype checks */
	if err := cmd.Header.VerifyStructure("http://schemas.microsoft.com/windows/pki/2009/01/enrollment/RST/wstep", true); err != nil {
		return err
	}

	if len(cmd.Body.BinarySecurityToken.Value) == 0 {
		return errors.New("invalid body binary security token")
	}

	// TODO: More stuff here

	/* Verify Context: Expensive checks against the server's DB */
	if err := cmd.Header.VerifyContext(config); err != nil {
		return err
	}

	// TODO: More stuff here

	return nil
}
