package certificate

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func DecodeStringToCertificate(raw string) (*x509.Certificate, error) {
	bytes, err := DecodeStringToASN1(raw)
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(bytes)
}

// Try to decode the PEM certificate in 3 step:
// * The certificate have his header and he is parse with the right length per line
// * The certificate have NOT his header but he is parse with the right length per line
// * The certificate have NOT his header and he is on one line
func DecodeStringToASN1(raw string) ([]byte, error) {
	block, _ := pem.Decode([]byte(raw))
	if block != nil {
		// has a valid PEM encoding
		return block.Bytes, nil
	}

	block, _ = pem.Decode([]byte(fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%v\n-----END CERTIFICATE-----", raw)))
	if block != nil {
		// was missing begin/end separators
		return block.Bytes, nil
	}

	bytes, err := decodeBase64([]byte(raw))
	if err == nil {
		// was a 1 line
		return bytes, nil
	}

	return nil, fmt.Errorf("invalid cerfificate %v", raw)
}

func decodeBase64(data []byte) ([]byte, error) {
	resultData := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(resultData, data)
	if err != nil {
		return nil, err
	}

	resultData = resultData[:n]

	return resultData, nil
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
