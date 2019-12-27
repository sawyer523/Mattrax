package boltdb

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"log"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mattrax/Mattrax/internal/types"
	"github.com/pkg/errors"
)

// certificatesBucket stores the name of the boltdb bucket that the certificates are stored in
var certificatesBucket = []byte("certificates")

var identityCertificateKey = []byte("identityCertificate")

var identityPrivateKeyKey = []byte("identityPrivateKey")

// CertificateService contains the implemented functionality for certificate management
type CertificateService struct {
	db *bolt.DB
}

func (cs CertificateService) GetIdentityRaw() ([]byte, []byte, error) {
	var certificateDer []byte
	var privateKeyDer []byte
	err := cs.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(certificatesBucket)
		if bucket == nil {
			return errors.New("error in CertificateService.GetIdentity: certificates bucket does not exist")
		}

		certificateDer = bucket.Get(identityCertificateKey)
		privateKeyDer = bucket.Get(identityPrivateKeyKey)
		if certificateDer == nil || privateKeyDer == nil {
			return errors.New("error CertificateService.GetIdentity: Mattrax Identity does not exist")
		}

		return nil
	})

	return certificateDer, privateKeyDer, err
}

func (cs CertificateService) GetIdentity() (*x509.Certificate, *rsa.PrivateKey, error) {
	certificateDer, privateKeyDer, err := cs.GetIdentityRaw()
	if err != nil {
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(certificateDer)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error CertificateService.GetIdentity: unable to decode the identity certificate")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyDer)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error CertificateService.GetIdentity: unable to decode the identity private key")
	}

	return cert, privateKey, nil
}

// rsaPublicKey reflects the ASN.1 structure of a PKCS#1 public key.
type rsaPublicKey struct {
	N *big.Int
	E int
}

// TODO: Cache identityCertificate, move this to Windows package
func (cs CertificateService) SignWSTEPRequest(binarySecurityToken string) ([]byte, *x509.Certificate, error) {
	identityCertificate, identityPrivKey, err := cs.GetIdentity()
	if err != nil {
		return nil, nil, err
	}

	// Decode Base64
	csrRaw, err := base64.StdEncoding.DecodeString(binarySecurityToken)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error CertificateService.SignWSTEPRequest: unable to decode binarySecurityToken")
	}

	// // Decode and verify CSR
	csr, err := x509.ParseCertificateRequest(csrRaw)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error CertificateService.SignWSTEPRequest: unable to parse certificate signing request")
	}
	if err = csr.CheckSignature(); err != nil {
		return nil, nil, errors.Wrap(err, "error CertificateService.SignWSTEPRequest: invalid certificate signing request signature")
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
		Issuer:             identityCertificate.Issuer,
		Subject:            csr.Subject,
		NotBefore:          NotBefore1,
		NotAfter:           NotBefore1.Add(365 * 24 * time.Hour),           // TODO: Configurable + Working renewal
		KeyUsage:           x509.KeyUsageDigitalSignature,                  // TODO: What does it do
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, // TODO: What does it do
	}

	clientCRTRaw, err := x509.CreateCertificate(rand.Reader, clientCertificate, identityCertificate, csr.PublicKey, identityPrivKey)
	if err != nil {
		return nil, nil, err // TODO: Wrap Error
	}

	return clientCRTRaw, clientCertificate, nil
}

// NewCertificateService creates and initialises a new CertificateService from a DB connection
func NewCertificateService(db *bolt.DB, certConfig types.IdentityCertificateConfig) (CertificateService, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(certificatesBucket)
		if err != nil {
			return errors.Wrap(err, "error NewCertificateService: error creating bucket")
		}

		identityCertificate := bucket.Get(identityCertificateKey)
		identityPrivateKey := bucket.Get(identityPrivateKeyKey)
		if identityCertificate == nil || identityPrivateKey == nil {
			log.Println("Creating new Mattrax Identity...")

			// Generate Private Key
			privateKey, err := rsa.GenerateKey(rand.Reader, certConfig.KeyLength)
			if err != nil {
				return errors.Wrap(err, "error NewCertificateService: generating identity private key")
			}

			// Generate SubjectKeyId, which is a 160-bit SHA-1 hash of the value of the public key
			publicKeyBytes, err := asn1.Marshal(rsaPublicKey{
				N: privateKey.PublicKey.N,
				E: privateKey.PublicKey.E,
			})
			if err != nil {
				return errors.Wrap(err, "error NewCertificateService: marshalling identity rsaPublicKey")
			}

			// TODO: Merge these two lines into one
			subjectKeyIDRaw := sha1.Sum(publicKeyBytes)
			subjectKeyID := subjectKeyIDRaw[:]

			// Create certificate template
			NotBefore := time.Now().Add(time.Duration(mathrand.Int31n(120)) * -time.Minute) // This randomises the creation time for security
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

			// Generate certificate
			certificateDer, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &privateKey.PublicKey, privateKey)
			if err != nil {
				return errors.Wrap(err, "error NewCertificateService: generating identity certificate")
			}

			// Store
			if err := bucket.Put(identityCertificateKey, certificateDer); err != nil {
				return err
			}
			if err := bucket.Put(identityPrivateKeyKey, x509.MarshalPKCS1PrivateKey(privateKey)); err != nil {
				return err
			}

			log.Println("Created Mattrax Identity")

			return nil
		}

		return err
	})

	return CertificateService{
		db,
	}, err
}

// TODO: Audit this entire file to check it follows MDM spec and impliments all security measures
