package enrollprovision

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"strings"
	"time"

	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/mattrax/Mattrax/internal/certificates"
	"github.com/mattrax/Mattrax/internal/devices"
)

// SignClientCertificate uses the Mattrax Identity CA to sign the CSR contained inside the Binary Security Token
// It also updates the device object to contain details about the certificate.
func SignClientCertificate(server *mattrax.Server, device devices.Device, binarySecurityToken string, identityCertificate certificates.Identity) ([]byte, error) {
	device.IdentityCertificate.NotBefore = time.Now().Add(time.Duration(mathrand.Int31n(120)) * -time.Minute) // This randomises the creation time a bit for added security (Recommended by x509 certificate signing not the MDM spec)
	device.IdentityCertificate.NotAfter = device.IdentityCertificate.NotBefore.Add(365 * 24 * time.Hour)
	if device.Windows.EnrollmentType == "Device" {
		device.IdentityCertificate.Subject.CommonName = device.Windows.DeviceID
	} else {
		device.IdentityCertificate.Subject.CommonName = device.EnrolledBy.Email
	}

	certificateSigningRequestDer, err := base64.StdEncoding.DecodeString(binarySecurityToken)
	if err != nil {
		return nil, err
	}

	certificateSigningRequest, err := x509.ParseCertificateRequest(certificateSigningRequestDer)
	if err != nil {
		return nil, err
	} else if err = certificateSigningRequest.CheckSignature(); err != nil {
		return nil, err
	}

	clientCertificate := &x509.Certificate{
		// TODO: Verify against other device certs
		Signature:          certificateSigningRequest.Signature,
		SignatureAlgorithm: certificateSigningRequest.SignatureAlgorithm,
		PublicKeyAlgorithm: certificateSigningRequest.PublicKeyAlgorithm,
		PublicKey:          certificateSigningRequest.PublicKey,
		SerialNumber:       big.NewInt(2), // TODO: ??? Increasing ???
		Issuer:             identityCertificate.Cert.Issuer,
		Subject:            device.IdentityCertificate.Subject,
		NotBefore:          device.IdentityCertificate.NotBefore,
		NotAfter:           device.IdentityCertificate.NotAfter,
		KeyUsage:           x509.KeyUsageDigitalSignature,
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	clientCertificateDer, err := x509.CreateCertificate(rand.Reader, clientCertificate, identityCertificate.Cert, certificateSigningRequest.PublicKey, identityCertificate.Key)
	if err != nil {
		return nil, err
	}

	h := sha1.New()
	h.Write(clientCertificateDer)
	device.IdentityCertificate.Hash = strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil))) // TODO: Cleanup -> This line is probally messer than it needs to be

	return clientCertificateDer, nil
}
