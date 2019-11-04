package hashicorp

import (
	"fmt"
	"path"

	"github.com/hashicorp/vault/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// LogicalV1 is a wrapper around api.Logical specialized in dealing with KVV1
type LogicalV1 struct {
	Base *api.Logical
	// AbsoluteSecretPath in which the logical can operate in the vault fs
	MountPoint string
	SecretPath string
}

// NewLogicalV1 construct a logical
func NewLogicalV1(impl *api.Logical, mountPoint, secretPath string) *LogicalV1 {
	// Check and append the path
	return &LogicalV1{
		Base:       impl,
		MountPoint: mountPoint,
		SecretPath: secretPath,
	}
}

// Write a secret inside the vault
func (l *LogicalV1) Write(subpath, value string) error {
	_, err := l.Base.Write(
		path.Join(l.MountPoint, l.SecretPath, subpath),
		map[string]interface{}{"value": value},
	)
	if err != nil {
		return errors.ConnectionError(err.Error()).SetComponent(component)
	}

	return nil
}

// Read a secret from the vault
func (l *LogicalV1) Read(subpath string) (value string, ok bool, err error) {
	res, err := l.Base.Read(path.Join(l.MountPoint, l.SecretPath, subpath))
	if err != nil {
		return "", false, errors.ConnectionError(err.Error()).SetComponent(component)
	}

	// When the secret is missing the client returns nil, nil.
	// We catch it here
	if res == nil {
		return "", false, nil
	}
	value = res.Data["value"].(string)

	return value, true, nil
}

// Delete remove the key from the vault
func (l *LogicalV1) Delete(subpath string) error {
	_, err := l.Base.Delete(path.Join(l.MountPoint, l.SecretPath, subpath))
	if err != nil {
		return errors.ConnectionError(err.Error()).SetComponent(component)
	}

	return nil
}

// List retrieve all the keys availables in the vault
func (l *LogicalV1) List(subpath string) ([]string, error) {
	res, err := l.Base.List(path.Join(l.MountPoint, l.SecretPath, subpath))
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
