package generic

import (
	mattrax "github.com/mattrax/Mattrax/internal"
)

type MdeWSSESecurity struct { // TODO: Verify type
	Username            string `xml:"wsse:UsernameToken>wsse:Username"`
	Password            string `xml:"wsse:UsernameToken>wsse:Password"`
	BinarySecurityToken string `xml:"wsse:BinarySecurityToken"`
}

type Header struct {
	Action    string `xml:"a:Action"`
	MessageID string `xml:"a:MessageID"`
	ReplyTo   struct {
		Address string `xml:"a:Address"`
	} `xml:"a:ReplyTo"`
	To           string          `xml:"a:To"`
	WSSESecurity MdeWSSESecurity `xml:"wsse:Security,omitempty"`
}

// TODO: Deprecated. Remove once unused.
func (header Header) VerifyStructure(action string, verifyWSSESecurity bool) error {
	return nil
}

// TODO: Deprecated. Remove once unused.
func (header Header) VerifyContext(config mattrax.Config) error {
	return nil
}
