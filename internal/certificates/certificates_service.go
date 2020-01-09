package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Service contains the code for safely (using a Mutex) getting and updating certificates.
type Service struct {
	certificates Certificates
	mutex        *sync.Mutex // Mutex is used to ensures exclusive access to the certificates
	store        Store
}

// Get returns the loaded settings.
func (s *Service) Get() Certificates {
	s.mutex.Lock()
	settings := s.certificates
	s.mutex.Unlock()

	return settings
}

// GenerateIdentity creates a new identity certificate and key pair
func (s *Service) GenerateIdentity(subject pkix.Name) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return errors.Wrap(err, "error generating identity private key")
	}

	publicKeyBytes, err := asn1.Marshal(RsaPublicKey{
		N: privateKey.PublicKey.N,
		E: privateKey.PublicKey.E,
	})
	if err != nil {
		return errors.Wrap(err, "error marshalling identity rsaPublicKey")
	}

	subjectKeyIDRaw := sha1.Sum(publicKeyBytes)
	subjectKeyID := subjectKeyIDRaw[:]

	NotBefore := time.Now().Add(time.Duration(mathrand.Int31n(120)) * -time.Minute) // This randomises the creation time for added security
	certificate := &x509.Certificate{
		SerialNumber:                big.NewInt(1), // TODO: What does it do. Should it be increasing
		Subject:                     subject,
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
		return errors.Wrap(err, "error generating identity certificate")
	}

	// TODO: Replace with sha1_hash := fmt.Sprintf("%x", sha1.Sum(form_value))
	sha1Hasher := sha1.New()
	sha1Hasher.Write(certificateDer)
	certificateHash := strings.ToUpper(fmt.Sprintf("%x", sha1Hasher.Sum(nil))) // TODO: Cleanup

	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)

	s.mutex.Lock()
	previousCertificates := s.certificates
	s.certificates.Identity = Identity{
		Cert:      certificate,
		CertRaw:   certificateDer,
		CertHash:  certificateHash,
		Key:       privateKey,
		KeyRaw:    privateKeyDer,
		Subject:   certificate.Subject,
		NotBefore: certificate.NotBefore,
		NotAfter:  certificate.NotAfter,
	}

	if err := s.store.Save(s.certificates); err != nil {
		s.certificates = previousCertificates
		s.mutex.Unlock()
		log.Error().Err(err).Msg("error saving new Mattrax identity")
		return errors.New("internal error saving certificates. intiial certificates were restored")

	}

	s.mutex.Unlock()
	log.Info().Str("CommonName", s.certificates.Identity.Subject.CommonName).Time("Expires", s.certificates.Identity.NotAfter).Msg("Generated new Mattrax identity...")
	return nil
}

// Store is a place where certificates are stored.
type Store interface {
	Save(Certificates) error
	Retrieve() (Certificates, error)
}

// NewService initialises and returns a new CertificateService
func NewService(store Store) (*Service, error) {
	certificates, err := store.Retrieve()
	if err != nil {
		return nil, err
	}

	return &Service{
		certificates: certificates,
		mutex:        &sync.Mutex{},
		store:        store,
	}, nil
}
