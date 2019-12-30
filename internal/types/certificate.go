package types

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"

	"github.com/pkg/errors"
)

// IdentityCertificateConfig contains the Identity certificate and its associated keys configurable values
type IdentityCertificateConfig struct {
	KeyLength int
	Subject   pkix.Name
}

// RsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type RsaPublicKey struct {
	N *big.Int
	E int
}

// ErrIdentityNotFound is returned when the identity certificate or private key don't exist
var ErrIdentityNotFound = errors.New("error: identity certificate or private key not found")

// CertificateStore is a storage mechanise capable of permanently storing certificates
type CertificateStore interface {
	RetrieveIdentity() (certRaw []byte, cert *x509.Certificate, keyRaw []byte, key *rsa.PrivateKey, err error)
	SaveIdentity(cert []byte, key *rsa.PrivateKey) error
}

// Certificates stores the certificates
type Certificates struct {
	IdentityCertRaw []byte
	IdentityKeyRaw  []byte
	IdentityCert    *x509.Certificate
	IdentityKey     *rsa.PrivateKey
}
