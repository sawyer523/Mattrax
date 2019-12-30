package generic

import (
	"net/url"
	"strings"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/pkg/errors"
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

func (header Header) VerifyStructure(action string, verifyWSSESecurity bool) error {
	if header.Action == "" {
		return errors.New("empty Action")
	} else if header.Action != action {
		return errors.New("invalid Action expected '" + action + "' but got '" + header.Action + "'")
	}

	if header.MessageID == "" {
		return errors.New("empty MessageID")
	} else if !strings.HasPrefix(header.MessageID, "urn:uuid:") {
		return errors.New("invalid MessageID prefix '" + header.MessageID + "'")
	} else if !validMessageID.MatchString(header.MessageID) {
		return errors.New("invalid characters in MessageID '" + header.MessageID + "'")
	}

	if header.To == "" {
		return errors.New("empty To")
	} else if _, err := url.ParseRequestURI(header.To); err != nil {
		return errors.New("invalid To '" + header.To + "'")
	}

	if verifyWSSESecurity {
		if header.WSSESecurity.BinarySecurityToken != "" { // Federated Authentication
			if !validBinarySecurityToken.MatchString(header.WSSESecurity.BinarySecurityToken) {
				return errors.New("invalid binary security token")
			}
		} else if header.WSSESecurity.Username != "" { // On Premise Authentication
			if !types.ValidEmail.MatchString(header.WSSESecurity.Username) {
				// Note: the incorrect username is not displayed as it could have been a user accidentally typing thier password
				return errors.New("invalid email address")
			}

			if header.WSSESecurity.Password == "" {
				return errors.New("invalid empty password")
			}
		} else {
			return errors.New("no supported client authentication details received")
		}
	}

	return nil
}

func (header Header) VerifyContext(config mattrax.Config) error {
	// Verify valid To address
	if toAddrRaw, err := url.Parse(header.To); err != nil {
		_ = toAddrRaw // TEMP
		// This should NEVER be called because the url is verified in VerifyStructure
		return errors.New("invalid To '" + header.To + "'")
	}

	/* TOOD: Fix for new Settings:      else if !(strings.ToLower(toAddrRaw.Hostname()) == strings.ToLower(config.Domain) || strings.ToLower(toAddrRaw.Hostname()) == strings.ToLower(config.WindowsDiscoveryDomain)) {
		return errors.New("this requested server ('" + strings.ToLower(toAddrRaw.Hostname()) + "') isn't this server ('" + config.WindowsDiscoveryDomain + "' or '" + config.PrimaryDomain + "')")
	} */

	return nil
}
