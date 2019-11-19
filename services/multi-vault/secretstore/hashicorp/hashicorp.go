package hashicorp

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/time/rate"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// HashiCorp is a wrapper around an HashiCorp object
type Hashicorp struct {
	// Client is the HashiCorp implementation of the vault client
	Client *api.Client
	// Logical is the query implementation of the vault
	Logical Logical
}

// Logical is a wrapper around api.Logical
type Logical interface {
	// Read value stored at given subpath
	Read(subpath string) (value string, ok bool, err error)
	// Write a value in the vault at given subpath
	Write(subpath, value string) error
	// List retrieve all the keys availables in the vault
	List(subpath string) ([]string, error)
	// Delete remove the key from the vault
	Delete(subpath string) error
}

// NewVaultClient returns a default object
func NewVaultClient(config *Config) (*Hashicorp, error) {
	// Instantiate an api.Client from local config object
	vaultConfig := ToVaultConfig(config)
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, errors.ConnectionError("Connection error: %v", err).
			SetComponent(component)
	}

	hashicorp := &Hashicorp{
		Client: client,
	}

	switch config.KVVersion {
	default:
		hashicorp.Logical = NewLogicalV2(client.Logical(), config.MountPoint, config.SecretPath)
	case "v1":
		hashicorp.Logical = NewLogicalV1(client.Logical(), config.MountPoint, config.SecretPath)
	}

	return hashicorp, nil
}

// ToVaultConfig extracts an api.Config object from self
func ToVaultConfig(c *Config) *api.Config {
	// Create Vault Configuration
	config := api.DefaultConfig()
	config.Address = c.Address
	config.HttpClient = cleanhttp.DefaultClient()
	config.HttpClient.Timeout = time.Second * 60

	// Create Transport
	transport := config.HttpClient.Transport.(*http.Transport)
	transport.TLSHandshakeTimeout = 10 * time.Second
	transport.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if err := http2.ConfigureTransport(transport); err != nil {
		config.Error = err
		return config
	}

	// Configure TLS
	tlsConfig := &api.TLSConfig{
		CACert:        c.CACert,
		CAPath:        c.CAPath,
		ClientCert:    c.ClientCert,
		ClientKey:     c.ClientKey,
		TLSServerName: c.TLSServerName,
		Insecure:      c.SkipVerify,
	}

	_ = config.ConfigureTLS(tlsConfig)

	config.Limiter = rate.NewLimiter(rate.Limit(c.RateLimit), c.BurstLimit)
	config.MaxRetries = c.MaxRetries
	config.Timeout = c.ClientTimeout

	// Ensure redirects are not automatically followed
	// Note that this is sane for the API client as it has its own
	// redirect handling logic (and thus also for command/meta),
	// but in e.g. http_test actual redirect handling is necessary
	config.HttpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		// Returning this value causes the Go net library to not close the
		// response body and to nil out the error. Otherwise retry clients may
		// try three times on every redirect because it sees an error from this
		// function (to prevent redirects) passing through to it.
		return http.ErrUseLastResponse
	}

	config.Backoff = retryablehttp.LinearJitterBackoff

	return config
}

// SetTokenFromConfig set the initial client token
func (v *Hashicorp) SetTokenFromConfig(c *Config) error {
	encoded, err := ioutil.ReadFile(c.TokenFilePath)
	if err != nil {
		log.Warningf("Token file path could not be found: %v", err.Error())
		return err
	}
	// Immediately delete the file after it was read
	_ = os.Remove(c.TokenFilePath)

	decoded := strings.TrimSuffix(string(encoded), "\n") // Remove the newline if it exists
	decoded = strings.TrimSuffix(decoded, "\r")          // This one is for windows compatibility
	v.Client.SetToken(decoded)

	return nil
}

// Auth is a shortcut to get the Auth object of the client
func (v *Hashicorp) Auth() *api.Auth {
	return v.Client.Auth()
}
