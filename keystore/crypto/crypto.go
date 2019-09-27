package crypto

import (
	"crypto/ecdsa"
)

// DSA is a digital signature algorithm. It is meant to be agnostic of ethereum
// but must be compatible
type DSA interface {
	Sign(digest []byte, priv *ecdsa.PrivateKey) ([]byte, error)
}
