package certificates

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

// RsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type RsaPublicKey struct {
	N *big.Int
	E int
}

// Certificates contains all the certificates for the server. It has both their raw and processed values.
type Certificates struct {
	Identity Identity
}

// Identity contains the certificates related to identifying and validating MDM clients
type Identity struct {
	Cert      *x509.Certificate `graphql:"-"`
	CertRaw   []byte
	CertHash  string          // SHA-1 hash of IdentityCertRaw
	Key       *rsa.PrivateKey `graphql:"-"`
	KeyRaw    []byte          `graphql:"-"`
	Subject   pkix.Name       `graphql:"-"`
	NotBefore time.Time
	NotAfter  time.Time
}
