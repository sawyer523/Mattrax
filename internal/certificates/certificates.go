package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/mattrax/Mattrax/internal/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Initialise loads the Identity certificate & private key and if not found generates a new certificate & private key
func Initialise(store types.CertificateStore, certConfig types.IdentityCertificateConfig) (types.Certificates, error) {
	certRaw, cert, keyRaw, key, err := store.RetrieveIdentity()
	if err != nil && err != types.ErrIdentityNotFound {
		return types.Certificates{}, err
	}

	if cert == nil || key == nil {
		log.Info().Msg("Generating new Mattrax identity...")

		privateKey, err := rsa.GenerateKey(rand.Reader, certConfig.KeyLength)
		if err != nil {
			return types.Certificates{}, errors.Wrap(err, "error generating identity private key")
		}

		// Generate SubjectKeyId, which is a 160-bit SHA-1 hash of the value of the public key
		publicKeyBytes, err := asn1.Marshal(types.RsaPublicKey{
			N: privateKey.PublicKey.N,
			E: privateKey.PublicKey.E,
		})
		if err != nil {
			return types.Certificates{}, errors.Wrap(err, "error marshalling identity rsaPublicKey")
		}

		subjectKeyIDRaw := sha1.Sum(publicKeyBytes)
		subjectKeyID := subjectKeyIDRaw[:]

		NotBefore := time.Now().Add(time.Duration(mathrand.Int31n(120)) * -time.Minute) // This randomises the creation time for added security
		certificate := &x509.Certificate{
			SerialNumber:                big.NewInt(1), // TODO: What does it do
			Subject:                     certConfig.Subject,
			NotBefore:                   NotBefore,
			NotAfter:                    NotBefore.Add(365 * 24 * time.Hour),
			KeyUsage:                    x509.KeyUsageCertSign | x509.KeyUsageCRLSign, // TODO: Are they required
			ExtKeyUsage:                 nil,                                          // TODO: What does it do
			UnknownExtKeyUsage:          nil,                                          // TODO: What does it do
			BasicConstraintsValid:       true,                                         // TODO: What does it do
			IsCA:                        true,
			MaxPathLen:                  0, // TODO: What does it do
			SubjectKeyId:                subjectKeyID,
			DNSNames:                    nil,
			PermittedDNSDomainsCritical: false, // TODO: What does it do
			PermittedDNSDomains:         nil,   // TODO: What does it do
		}

		certificateDer, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &privateKey.PublicKey, privateKey)
		if err != nil {
			return types.Certificates{}, errors.Wrap(err, "error NewCertificateService: generating identity certificate")
		}

		store.SaveIdentity(certificateDer, privateKey)

		log.Info().Msg("Generated new Mattrax identity...")

		return types.Certificates{
			IdentityCertRaw: certificateDer,
			IdentityCert:    certificate,
			IdentityKeyRaw:  x509.MarshalPKCS1PrivateKey(privateKey),
			IdentityKey:     privateKey,
		}, nil
	}

	return types.Certificates{
		IdentityCertRaw: certRaw,
		IdentityCert:    cert,
		IdentityKeyRaw:  keyRaw,
		IdentityKey:     key,
	}, nil
}
