package types

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
)

type IdentityCertificateConfig struct {
	KeyLength int
	Subject   pkix.Name
}

// CertificateService contains the implemented functionality for certificate management
type CertificateService interface {
	GetIdentityRaw() (certificateDer []byte, privateKeyDer []byte, err error)
	GetIdentity() (*x509.Certificate, *rsa.PrivateKey, error)
	SignWSTEPRequest(binarySecurityToken string) (signedCertDer []byte, cert *x509.Certificate, err error)
}
