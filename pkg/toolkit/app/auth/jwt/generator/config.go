package generator

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/jwt"
	"github.com/consensys/orchestrate/pkg/toolkit/tls/certificate"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(PrivateKeyViperKey, privateKeyDefault)
	_ = viper.BindEnv(PrivateKeyViperKey, privateKeyEnv)
}

// Provision trusted certificate of the authentication service (base64 encoded)
const (
	privateKeyFlag     = "auth-jwt-key"
	PrivateKeyViperKey = "auth.jwt.key"
	privateKeyDefault  = ""
	privateKeyEnv      = "AUTH_JWT_KEY"
)

func PrivateKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path to private key to sign generated JWT tokens
Environment variable: %q`, privateKeyEnv)
	f.String(privateKeyFlag, privateKeyDefault, desc)
	_ = viper.BindPFlag(PrivateKeyViperKey, f.Lookup(privateKeyFlag))
}

type Config struct {
	KeyPair              *certificate.KeyPair
	OrchestrateClaimPath string
}

func NewConfig(vipr *viper.Viper) (*Config, error) {
	var certB []byte
	var keyB []byte
	if certFile := vipr.GetString(jwt.CertificateFileViperKey); certFile != "" {
		_, err := os.Stat(certFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read cert file. %s", err.Error())
		}

		certB, err = ioutil.ReadFile(certFile)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("empty cert file")
	}

	if keyFile := vipr.GetString(PrivateKeyViperKey); keyFile != "" {
		_, err := os.Stat(keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file. %s", err.Error())
		}

		keyB, err = ioutil.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("empty key file")
	}

	return &Config{
		KeyPair: &certificate.KeyPair{
			Cert: certB,
			Key:  keyB,
		},
		OrchestrateClaimPath: vipr.GetString(jwt.OrchestrateClaimPathViperKey),
	}, nil
}
