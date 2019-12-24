package wtypes

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/mattrax/Mattrax/pkg/xml"
)

// MdeGetPoliciesWSSESecurity contains the authentication details
type MdeGetPoliciesWSSESecurity struct {
	Username string `xml:"wsse:UsernameToken>wsse:Username"`
	Password string `xml:"wsse:UsernameToken>wsse:Password"`
}

// MdeGetPoliciesRequest contains the payload sent by the client to finish authentication and retieve policies
type MdeGetPoliciesRequest struct {
	XMLName xml.Name `xml:"s:Envelope"`
	Header  struct {
		Action    string `xml:"a:Action"`
		MessageID string `xml:"a:MessageID"`
		ReplyTo   struct {
			Address string `xml:"a:Address"`
		} `xml:"a:ReplyTo"`
		To           string                     `xml:"a:To"`
		WSSESecurity MdeGetPoliciesWSSESecurity `xml:"wsse:Security"`
	} `xml:"s:Header"`
}

// validPassword is a regex used to verify a Password is valid
var validPassword = regexp.MustCompile(`^[a-zA-Z0-9:\-@ !#$^&*().,?]+$`) // TODO: Use through user stuff, move package form here

// VerifyStructure checks that the cmd contains a valid structure.
// This is done before more resource intensive checks to reduce the DoS vector.
func (cmd MdeGetPoliciesRequest) VerifyStructure() error {
	if cmd.Header.Action == "" {
		return errors.New("invalid MdeGetPoliciesRequest: empty Action")
	} else if cmd.Header.Action != "http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy/IPolicy/GetPolicies" {
		return errors.New("invalid MdeGetPoliciesRequest: invalid Action '" + cmd.Header.Action + "'")
	}

	if cmd.Header.MessageID == "" {
		return errors.New("invalid MdeGetPoliciesRequest: empty MessageID")
	} else if !strings.HasPrefix(cmd.Header.MessageID, "urn:uuid:") {
		return errors.New("invalid MdeGetPoliciesRequest: invalid MessageID prefix '" + cmd.Header.MessageID + "'")
	} else if !validMessageID.MatchString(cmd.Header.MessageID) {
		return errors.New("invalid MdeGetPoliciesRequest: invalid characters in MessageID '" + cmd.Header.MessageID + "'")
	}

	if cmd.Header.To == "" {
		return errors.New("invalid MdeGetPoliciesRequest: empty To")
	} else if _, err := url.ParseRequestURI(cmd.Header.To); err != nil {
		return errors.New("invalid MdeGetPoliciesRequest: invalid To '" + cmd.Header.To + "'")
	}

	if cmd.Header.WSSESecurity.Username == "" {
		return errors.New("invalid MdeGetPoliciesRequest: empty email address")
	} else if !validEmail.MatchString(cmd.Header.WSSESecurity.Username) {
		// Note: the incorrect username is not displayed as it could have been a user accidentally typing thier password
		return errors.New("invalid MdeGetPoliciesRequest: invalid email address")
	}

	if cmd.Header.WSSESecurity.Password == "" {
		return errors.New("invalid MdeGetPoliciesRequest: empty password")
	} else if !validPassword.MatchString(cmd.Header.WSSESecurity.Password) {
		// Note: the incorrect password is not displayed for what should be obvious reasons
		return errors.New("invalid MdeGetPoliciesRequest: invalid password")
	}

	return nil
}

// VerifyContext checks that the cmd against the expected values.
// This operation is more expensive so should always be done after VerifyStructure.
func (cmd MdeGetPoliciesRequest) VerifyContext(config mattrax.Config, userService types.UserService) error {
	// Verify valid To address
	if toAddrRaw, err := url.Parse(cmd.Header.To); err != nil {
		// This should NEVER be called because the url is verified in VerifyStructure
		return errors.New("invalid MdeGetPoliciesRequest: invalid To '" + cmd.Header.To + "'")
	} else if strings.ToLower(toAddrRaw.Hostname()) != strings.ToLower(config.PrimaryDomain) {
		return errors.New("invalid MdeGetPoliciesRequest: this requested server ('" + cmd.Header.To + "') isn't this server ('" + config.WindowsDiscoveryDomain + "')")
	}

	// Verify user login
	loggedIn, err := userService.VerifyLogin(cmd.Header.WSSESecurity.Username, cmd.Header.WSSESecurity.Password)
	if err != nil {
		return err
	}

	if loggedIn {
		return nil
	}

	return errors.New("invalid MdeGetPoliciesRequest: the users login is incorrect")
}

type MdeCA struct {
	// TODO: file:///Users/oscar.beaumont/Downloads/[MS-XCEP].pdf#page=22&zoom=100,92,322
}

type MdeCACollection struct {
	Nil   bool    `xml:"xsi:nil,attr"`
	Value []MdeCA `xml:",innerxml"`
}

// MdeAttributesCertificateValidity contains information about the expected validity of an issued certificate
// Reference: MS-XCEP 3.1.4.1.3.8 CertificateValidity
// Fields must be a positive nonzero number.
type MdeAttributesCertificateValidity struct {
	ValidityPeriodSeconds int `xml:"validityPeriodSeconds"`
	RenewalPeriodSeconds  int `xml:"renewalPeriodSeconds"`
}

// MdeEnrollmentPermission is used to convey the permissions for the associated parent object
// Reference: 3.1.4.1.3.11 EnrollmentPermission
type MdeEnrollmentPermission struct {
	Enroll     bool `xml:"enroll"`
	AutoEnroll bool `xml:"autoEnroll"`
}

type MdeKeySpec struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdeKeyUsageProperty struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdePermissions struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdeAlgorithmOIDReference struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdeCryptoProviders struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

// MdePrivateKeyAttributes contains the attributes for the private key that will be associated with any certificate request
// Reference 3.1.4.1.3.20 PrivateKeyAttributes
type MdePrivateKeyAttributes struct {
	// MinimalKeyLength must be a positive nonzero number
	MinimalKeyLength      int
	KeySpec               MdeKeySpec               `xml:"keySpec"`
	KeyUsageProperty      MdeKeyUsageProperty      `xml:"keyUsageProperty"`
	Permissions           MdePermissions           `xml:"permissions"`
	AlgorithmOIDReference MdeAlgorithmOIDReference `xml:"algorithmOIDReference"`
	CryptoProviders       MdeCryptoProviders       `xml:"cryptoProviders"`
}

// MdeRevision identifies the version of the associated MdePolicy
// Reference 3.1.4.1.3.24 Revision
type MdeRevision struct {
	// MajorRevision must be a positive nonzero integer
	MajorRevision int `xml:"majorRevision"`
	// MinorRevision must be an integer greater than or equal to 0
	MinorRevision int `xml:"minorRevision"`
}

type MdeSupersededPolicies struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdePrivateKeyFlags struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdeSubjectNameFlags struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdeEnrollmentFlags struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdeGeneralFlags struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdeRARequirements struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdeKeyArchivalAttributes struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdeExtensions struct {
	Nil bool `xml:"xsi:nil,attr"`
	// TODO
}

type MdeAttributes struct {
	// CommonName must be unique in current GetPoliciesResponse
	CommonName                string                           `xml:"commonName"`
	PolicySchema              int                              `xml:"policySchema"`
	CertificateValidity       MdeAttributesCertificateValidity `xml:"certificateValidity"`
	EnrollmentPermission      MdeEnrollmentPermission          `xml:"permission"`
	PrivateKeyAttributes      MdePrivateKeyAttributes          `xml:"privateKeyAttributes"`
	Revision                  MdeRevision                      `xml:"revision"`
	SupersededPolicies        MdeSupersededPolicies            `xml:"supersededPolicies"`
	PrivateKeyFlags           MdePrivateKeyFlags               `xml:"privateKeyFlags"`
	SubjectNameFlags          MdeSubjectNameFlags              `xml:"subjectNameFlags"`
	EnrollmentFlags           MdeEnrollmentFlags               `xml:"enrollmentFlags"`
	GeneralFlags              MdeGeneralFlags                  `xml:"generalFlags"`
	HashAlgorithmOIDReference int                              `xml:"hashAlgorithmOIDReference"`
	RARequirements            MdeRARequirements                `xml:"rARequirements"`
	KeyArchivalAttributes     MdeKeyArchivalAttributes         `xml:"keyArchivalAttributes"`
	Extensions                MdeExtensions                    `xml:"extensions"`
}

type MdePolicy struct {
	XMLName      xml.Name        `xml:"policy"`
	OIDReference int             `xml:"policyOIDReference"` // Unique policy "id"
	CAs          MdeCACollection `xml:"cAs,omitempty"`
	Attributes   MdeAttributes   `xml:"attributes"`
}

// DiscoverResponse is the payload within the body of the response
type MdePoliciesResponse struct {
	PolicyID           string      `xml:"xcep:response>policyID"`
	PolicyFriendlyName string      `xml:"xcep:response>policyFriendlyName"`
	NextUpdateHours    int         `xml:"xcep:response>nextUpdateHours,omitempty"`
	PoliciesNotChanged bool        `xml:"xcep:response>policiesNotChanged,omitempty"`
	Policies           []MdePolicy `xml:"xcep:response>policies"`
	// TODO: Verify CAS + OIDS types are correct here
	CAS  string `xml:"xcep:cAs"`
	OIDS string `xml:"xcep:oIDs"`
}

// MdePolicyResponseBody is the body tag of the response
type MdePolicyResponseBody struct {
	NamespaceXSI     string              `xml:"xmlns:xsi,attr"`
	NamespaceXSD     string              `xml:"xmlns:xsd,attr"`
	PoliciesResponse MdePoliciesResponse `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy GetPoliciesResponse"`
}

// MdePolicyResponseEnvelope contains the payload sent from the server to the client telling the client of its enrollment policies and certificate templates
type MdePolicyResponseEnvelope struct {
	XMLName    xml.Name `xml:"s:Envelope"`
	NamespaceS string   `xml:"xmlns:s,attr"`
	NamespaceA string   `xml:"xmlns:a,attr"`
	NamespaceU string   `xml:"xmlns:u,attr"`

	HeaderAction MustUnderstand `xml:"s:Header>a:Action"`
	// HeaderActivityID string                `xml:"s:Header>a:ActivityID"` // TODO: Should this be included as the docs don't show it
	HeaderRelatesTo string                `xml:"s:Header>a:RelatesTo"`
	Body            MdePolicyResponseBody `xml:"s:Body"`
}
