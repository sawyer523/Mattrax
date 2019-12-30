package wstep

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/mattrax/Mattrax/internal/types"
	"github.com/pkg/errors"
)

func SignRequest(cert types.Certificates, binarySecurityToken string) ([]byte, *x509.Certificate, error) {
	// Decode Base64
	csrRaw, err := base64.StdEncoding.DecodeString(binarySecurityToken)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error unable to decode binarySecurityToken")
	}

	// // Decode and verify CSR
	csr, err := x509.ParseCertificateRequest(csrRaw)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error unable to parse certificate signing request")
	}
	if err = csr.CheckSignature(); err != nil {
		return nil, nil, errors.Wrap(err, "error invalid certificate signing request signature")
	}

	// // Create client identity certificate template
	// // TODO: Remove 1 from name
	NotBefore1 := time.Now().Add(time.Duration(mathrand.Int31n(120)) * -time.Minute) // This randomises the creation time for security
	clientCertificate := &x509.Certificate{
		Signature:          csr.Signature,
		SignatureAlgorithm: csr.SignatureAlgorithm,
		PublicKeyAlgorithm: csr.PublicKeyAlgorithm,
		PublicKey:          csr.PublicKey,
		SerialNumber:       big.NewInt(2), // TODO: What does it do, should it be increasing?
		Issuer:             cert.IdentityCert.Issuer,
		Subject:            csr.Subject,
		NotBefore:          NotBefore1,
		NotAfter:           NotBefore1.Add(365 * 24 * time.Hour),           // TODO: Configurable + Working renewal
		KeyUsage:           x509.KeyUsageDigitalSignature,                  // TODO: What does it do
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, // TODO: What does it do
	}

	clientCRTRaw, err := x509.CreateCertificate(rand.Reader, clientCertificate, cert.IdentityCert, csr.PublicKey, cert.IdentityKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error creating certificate")
	}

	return clientCRTRaw, clientCertificate, nil
}
