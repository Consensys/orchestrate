package hashicorp

import (
	"fmt"
	"path"

	"github.com/hashicorp/vault/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

// LogicalV2 is a wrapper around api.Logical specialized in dealing with KVV1
type LogicalV2 struct {
	Base       *api.Logical
	MountPoint string
	SecretPath string
}

// NewLogicalV2 construct a logical
func NewLogicalV2(impl *api.Logical, mountPoint, secretPath string) *LogicalV2 {
	// Check and append the path
	return &LogicalV2{
		Base:       impl,
		MountPoint: mountPoint,
		SecretPath: secretPath,
	}
}

// Read fetch the value from vault by key
func (l *LogicalV2) Read(subPath string) (value string, ok bool, err error) {
	// Read secret from Vault
	res, err := l.Base.Read(
		path.Join(l.MountPoint, "data", l.SecretPath, subPath),
	)
	if err != nil {
		return "", false, errors.ConnectionError(err.Error()).SetComponent(component)
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
func (l *LogicalV2) Write(subPath, value string) error {
	// Load secret to Vault
	_, err := l.Base.Write(
		path.Join(l.MountPoint, "data", l.SecretPath, subPath),
		map[string]interface{}{
			"data": map[string]interface{}{"value": value},
		},
	)

	if err != nil {
		return errors.ConnectionError(err.Error()).SetComponent(component)
	}

	return nil
}

// Delete remove the key from the vault
func (l *LogicalV2) Delete(subPath string) error {
	// Delete secret in Vault
	_, err := l.Base.Delete(
		path.Join(l.MountPoint, "metadata", l.SecretPath, subPath),
	)
	if err != nil {
		return errors.ConnectionError(err.Error()).SetComponent(component)
	}

	return nil
}

// List retrieve all the keys availables in the vault
func (l *LogicalV2) List(subPath string) ([]string, error) {

	res, err := l.Base.List(path.Join(l.MountPoint, "metadata", l.SecretPath))
	if err != nil {
		return nil, errors.ConnectionError(err.Error()).SetComponent(component)
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
