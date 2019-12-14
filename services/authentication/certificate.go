package authentication

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

const certificateType = "CERTIFICATE"

func decodeBase64(data []byte) ([]byte, error) {
	resultData := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(resultData, data)
	if err != nil {
		return nil, err
	}

	resultData = resultData[:n]

	return resultData, nil
}

// Try to decode the PEM certificate in 3 step:
// * The certificate have his header and he is parse with the right length per line
// * The certificate have NOT his header but he is parse with the right length per line
// * The certificate have NOT his header and he is on one line
func getCert() (*x509.Certificate, error) {
	cert := ""
	var block *pem.Block

	rawCert := viper.GetString(AuthServiceCertificateViperKey)
	// Parse PEM block
	if block, _ = pem.Decode([]byte(cert)); block == nil {

		cert = "-----BEGIN CERTIFICATE-----\n" + rawCert + "\n-----END CERTIFICATE-----"
		if block, _ = pem.Decode([]byte(cert)); block == nil {
			block = &pem.Block{
				Type:  certificateType,
				Bytes: []byte(rawCert),
			}

			var err error
			block.Bytes, err = decodeBase64(block.Bytes)
			if err != nil {
				return nil, err
			}

			pemEncode := pem.EncodeToMemory(block)

			if block, _ = pem.Decode(pemEncode); block == nil {
				return nil, jwt.ErrKeyMustBePEMEncoded
			}
		}
	}

	// Parse the certificate
	decodedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return decodedCert, nil
}

// Convert a PublicKey to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

// Extract and return the key to use to validate the ID / Access Token
// Implement the interface of Keyfunc from jwt package
func ValidatedKey(token *jwt.Token) (interface{}, error) {
	cert, err := getCert()
	if err != nil {
		return nil, errors.New("unable to retrieve trusted certificate")
	}

	var publicKey *rsa.PublicKey
	var ok bool
	if publicKey, ok = cert.PublicKey.(*rsa.PublicKey); !ok {
		return nil, jwt.ErrNotRSAPublicKey
	}

	return publicKey, nil
}
