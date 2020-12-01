package kvv2

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/hashicorp/vault/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/hashicorp"
)

// LogicalV2 is a wrapper around api.Logical specialized in dealing with KVV1
type Client struct {
	Base       *api.Logical
	MountPoint string
	SecretPath string
	Health     func() (*api.HealthResponse, error)
}

func NewClient(cfg *hashicorp.Config, secretPath string) (*Client, error) {
	vaultConfig := hashicorp.ToVaultConfig(cfg)
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		return nil, errors.ConnectionError("Connection error: %v", err)
	}

	encoded, err := ioutil.ReadFile(cfg.TokenFilePath)
	if err != nil {
		return nil, err
	}

	decoded := strings.TrimSuffix(string(encoded), "\n") // Remove the newline if it exists
	decoded = strings.TrimSuffix(decoded, "\r")          // This one is for windows compatibility
	client.SetToken(decoded)

	return &Client{
		Base:       client.Logical(),
		Health:     client.Sys().Health,
		MountPoint: cfg.MountPoint,
		SecretPath: secretPath,
	}, nil
}

// Read fetch the value from vault by key
func (l *Client) Read(subPath string) (value string, ok bool, err error) {
	// Read secret from Vault
	res, err := l.Base.Read(
		path.Join(l.MountPoint, "data", l.SecretPath, subPath),
	)
	if err != nil {
		return "", false, errors.ConnectionError(err.Error())
	}

	// When the secret does not exist the client returns nil, nil.
	// When the secret is deleted the keystore returns an empty data.
	// We catch it here
	if res == nil || res.Data["data"] == nil {
		return "", false, nil
	}
	value = res.Data["data"].(map[string]interface{})["value"].(string)

	return value, true, nil
}

// Write the Secret value stored in the vault
func (l *Client) Write(subPath, value string) error {
	// Load secret to Vault
	_, err := l.Base.Write(
		path.Join(l.MountPoint, "data", l.SecretPath, subPath),
		map[string]interface{}{
			"data": map[string]interface{}{"value": value},
		},
	)

	if err != nil {
		return errors.ConnectionError(err.Error())
	}

	return nil
}

// Delete remove the key from the vault
func (l *Client) Delete(subPath string) error {
	// Delete secret in Vault
	_, err := l.Base.Delete(
		path.Join(l.MountPoint, "metadata", l.SecretPath, subPath),
	)
	if err != nil {
		return errors.ConnectionError(err.Error())
	}

	return nil
}

// List retrieve all the keys available in the vault
func (l *Client) List(subPath string) ([]string, error) {
	res, err := l.Base.List(path.Join(l.MountPoint, "metadata", l.SecretPath))
	if err != nil {
		return nil, errors.ConnectionError(err.Error())
	}

	if res == nil {
		return []string{}, nil
	}

	secrets := res.Data["keys"].([]interface{})
	rv := make([]string, len(secrets))
	for i, elem := range secrets {
		rv[i] = fmt.Sprintf("%v", elem)
	}

	return rv, nil
}
