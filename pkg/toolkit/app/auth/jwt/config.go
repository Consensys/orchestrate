package jwt

import (
	"context"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/toolkit/tls/certificate"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gopkg.in/square/go-jose.v2"
)

func init() {
	viper.SetDefault(CertificateViperKey, certificateDefault)
	_ = viper.BindEnv(CertificateViperKey, certificateEnv)
	viper.SetDefault(CertificateFileViperKey, certificateFileDefault)
	_ = viper.BindEnv(CertificateFileViperKey, certificateFileEnv)
	viper.SetDefault(CertificateIssuerURLViperKey, certificateIssuerURLDefault)
	_ = viper.BindEnv(CertificateIssuerURLViperKey, certificateIssuerURLEnv)
	viper.SetDefault(OrchestrateClaimPathViperKey, OrchestrateClaimPathDefault)
	_ = viper.BindEnv(OrchestrateClaimPathViperKey, OrchestrateClaimPathEnv)
}

func Flags(f *pflag.FlagSet) {
	certificateFlags(f)
	certificateFileFlags(f)
	certificateIssuerURLFlags(f)
	OrchestrateClaimPath(f)
}

// Provision trusted certificate of the authentication service (base64 encoded)
const (
	certificateFlag     = "auth-jwt-certificate"
	CertificateViperKey = "auth.jwt.certificate"
	certificateDefault  = ""
	certificateEnv      = "AUTH_JWT_CERTIFICATE"
)

// certificateFlag register flag for Authentication service certificate
func certificateFlags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`certificate of the authentication service encoded in base64 (DEPRECATED).
Environment variable: %q`, certificateEnv)
	f.String(certificateFlag, certificateDefault, desc)
	_ = viper.BindPFlag(CertificateViperKey, f.Lookup(certificateFlag))
}

// Provision trusted certificate of the authentication service (file path)
const (
	certificateFileFlag     = "auth-jwt-cert-file"
	CertificateFileViperKey = "auth.jwt.cert"
	certificateFileDefault  = ""
	certificateFileEnv      = "AUTH_JWT_CERT"
)

func certificateFileFlags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Path to certificate of the authentication service.
Environment variable: %q`, certificateFileEnv)
	f.String(certificateFileFlag, certificateFileDefault, desc)
	_ = viper.BindPFlag(CertificateFileViperKey, f.Lookup(certificateFileFlag))
}

const (
	certificateIssuerURLFlag     = "auth-jwt-issuer-url"
	CertificateIssuerURLViperKey = "auth.jwt.issuer-url"
	certificateIssuerURLDefault  = ""
	certificateIssuerURLEnv      = "AUTH_JWT_ISSUER_URL"
)

func certificateIssuerURLFlags(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`JWT issuer server domain (ie. https://orchestrate.eu.auth0.com/.well-known/jwks.json).
Environment variable: %q`, certificateIssuerURLEnv)
	f.String(certificateIssuerURLFlag, certificateIssuerURLDefault, desc)
	_ = viper.BindPFlag(CertificateIssuerURLViperKey, f.Lookup(certificateIssuerURLFlag))
}

// Provision tenant namespace to retrieve the tenant id in the OpenId or Access Token (JWT)
const (
	OrchestrateClaimPathFlag     = "auth-jwt-claims-namespace"
	OrchestrateClaimPathViperKey = "auth.jwt.claims.namespace"
	OrchestrateClaimPathDefault  = ""
	OrchestrateClaimPathEnv      = "AUTH_JWT_CLAIMS_NAMESPACE"
)

// OrchestrateClaimPath register flag for tenant namespace
func OrchestrateClaimPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Tenant Namespace to retrieve the tenant id in the Access Token (DEPRECATED).
Environment variable: %q`, OrchestrateClaimPathEnv)
	f.String(OrchestrateClaimPathFlag, OrchestrateClaimPathDefault, desc)
	_ = viper.BindPFlag(OrchestrateClaimPathViperKey, f.Lookup(OrchestrateClaimPathFlag))
}

type Config struct {
	Certificates         []*x509.Certificate
	OrchestrateClaimPath string
	SkipClaimsValidation bool
	ValidMethods         []string
}

func NewConfig(vipr *viper.Viper) (*Config, error) {
	cfg := &Config{
		OrchestrateClaimPath: vipr.GetString(OrchestrateClaimPathViperKey),
	}

	var err error
	if issuerURL := vipr.GetString(CertificateIssuerURLViperKey); issuerURL != "" {
		cfg.Certificates, err = oidcIssuerURLCert(issuerURL)
		if err != nil {
			return nil, fmt.Errorf("failed to read issuer cert. %s", err.Error())
		}
	} else if certFile := vipr.GetString(CertificateFileViperKey); certFile != "" {
		_, err = os.Stat(certFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read cert file. %s", err.Error())
		}

		certContent, err := ioutil.ReadFile(certFile)
		if err != nil {
			return nil, err
		}

		cert, err := decodeCertificate(certContent)
		if err != nil {
			return nil, err
		}
		cfg.Certificates = []*x509.Certificate{cert}
	} else if certStr := vipr.GetString(CertificateFileViperKey); certStr != "" {
		cert, err := decodeCertificate([]byte(certStr))
		if err != nil {
			return nil, err
		}

		cfg.Certificates = []*x509.Certificate{cert}
	}

	return cfg, nil
}

func decodeCertificate(content []byte) (*x509.Certificate, error) {
	bCert, err := certificate.Decode(content, "CERTIFICATE")
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(bCert[0])
}

func oidcIssuerURLCert(issuerServer string) ([]*x509.Certificate, error) {
	jwks, err := retrieveKeySet(context.Background(), http.DefaultClient, issuerServer)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve auth server jwks: %s", issuerServer)
	}

	if len(jwks.Keys) == 0 {
		return nil, fmt.Errorf("empty server jwks")
	}

	var certs []*x509.Certificate
	// nolint
	for _, kw := range jwks.Keys {
		certs = append(certs, kw.Certificates...)
	}

	return certs, nil
}

func retrieveKeySet(ctx context.Context, client *http.Client, authEndpoint string) (*jose.JSONWebKeySet, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", authEndpoint, nil)
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call to JWK server failed. %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to retrieve keys from JWK server")
	}

	keySet := &jose.JSONWebKeySet{}
	if err := json.UnmarshalBody(response.Body, keySet); err != nil {
		return nil, fmt.Errorf("failed to decode response body. %s", err.Error())
	}

	return keySet, nil
}
