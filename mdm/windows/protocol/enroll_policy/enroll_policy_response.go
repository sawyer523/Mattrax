package enrollpolicy

import (
	"github.com/mattrax/Mattrax/mdm/windows/soap"
	"github.com/mattrax/Mattrax/pkg/xml"
)

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

type Response struct {
	PolicyID           string      `xml:"xcep:response>policyID"`
	PolicyFriendlyName string      `xml:"xcep:response>policyFriendlyName"`
	NextUpdateHours    int         `xml:"xcep:response>nextUpdateHours,omitempty"`
	PoliciesNotChanged bool        `xml:"xcep:response>policiesNotChanged,omitempty"`
	Policies           []MdePolicy `xml:"xcep:response>policies"`
	// TODO: Verify CAS + OIDS types are correct here
	CAS  string `xml:"xcep:cAs"`
	OIDS string `xml:"xcep:oIDs"`
}

type ResponseBody struct {
	NamespaceXSI     string   `xml:"xmlns:xsi,attr"`
	NamespaceXSD     string   `xml:"xmlns:xsd,attr"`
	PoliciesResponse Response `xml:"http://schemas.microsoft.com/windows/pki/2009/01/enrollmentpolicy GetPoliciesResponse"`
}

type ResponseEnvelope struct {
	XMLName    xml.Name `xml:"s:Envelope"`
	NamespaceS string   `xml:"xmlns:s,attr"`
	NamespaceA string   `xml:"xmlns:a,attr"`
	NamespaceU string   `xml:"xmlns:u,attr"`

	HeaderAction soap.MustUnderstand `xml:"s:Header>a:Action"`
	// HeaderActivityID string                `xml:"s:Header>a:ActivityID"` // TODO: Should this be included as the docs don't show it
	HeaderRelatesTo string       `xml:"s:Header>a:RelatesTo"`
	Body            ResponseBody `xml:"s:Body"`
}
