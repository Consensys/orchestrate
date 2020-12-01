package kvv2

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/hashicorp"
)

type Config struct {
	*hashicorp.Config
	SecretPath string
}

func init() {
	viper.SetDefault(vaultMountPointViperKey, vaultMountPointDefault)
	_ = viper.BindEnv(vaultMountPointViperKey, vaultMountPointEnv)

	viper.SetDefault(vaultSecretPathViperKey, vaultSecretPathDefault)
	_ = viper.BindEnv(vaultSecretPathViperKey, vaultSecretPathEnv)

	viper.SetDefault(vaultTokenFilePathViperKey, vaultTokenFilePathDefault)
	_ = viper.BindEnv(vaultTokenFilePathViperKey, vaultTokenFilePathEnv)
}

const (
	vaultSecretPathViperKey = "vault.v2.secret.path"
	vaultSecretPathDefault  = "default"
	vaultSecretPathEnv      = "VAULT_V2_SECRET_PATH"
	vaultSecretPathFlag     = "vault-v2-secret-path"

	vaultMountPointViperKey = "vault.v2.mount.point"
	vaultMountPointDefault  = "secret"
	vaultMountPointFlag     = "vault-v2-mount-point"
	vaultMountPointEnv      = "VAULT_V2_MOUNT_POINT"

	vaultTokenFilePathViperKey = "vault.v2.token.file"
	vaultTokenFilePathDefault  = "/vault/token/.root"
	vaultTokenFilePathFlag     = "vault-v2-token-file"
	vaultTokenFilePathEnv      = "VAULT_V2_TOKEN_FILE"
)

func InitFlags(f *pflag.FlagSet) {
	hashicorp.InitFlags(f)
	vaultMountPoint(f)
	vaultTokenFilePath(f)
	vaultSecretPath(f)
}

func vaultMountPoint(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies the mount point used (v2 engine). Should not start with a //
Environment variable: %q `, vaultMountPointEnv)
	f.String(vaultMountPointFlag, vaultMountPointDefault, desc)
	_ = viper.BindPFlag(vaultMountPointViperKey, f.Lookup(vaultMountPointFlag))
}

func vaultTokenFilePath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Specifies the token file path (v2 engine).
Parameter ignored if the token has been passed by VAULT_TOKEN
Environment variable: %q `, vaultTokenFilePathEnv)
	f.String(vaultTokenFilePathFlag, vaultTokenFilePathDefault, desc)
	_ = viper.BindPFlag(vaultTokenFilePathViperKey, f.Lookup(vaultTokenFilePathFlag))
}

func vaultSecretPath(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret path (v2 engine).
Environment variable: %q`, vaultSecretPathEnv)
	f.String(vaultSecretPathFlag, vaultSecretPathDefault, desc)
	_ = viper.BindPFlag(vaultSecretPathViperKey, f.Lookup(vaultSecretPathFlag))
}

func ConfigFromViper() *Config {
	cfg := hashicorp.ConfigFromViper()
	cfg.MountPoint = viper.GetString(vaultMountPointViperKey)
	cfg.TokenFilePath = viper.GetString(vaultTokenFilePathViperKey)

	return &Config{
		Config:     cfg,
		SecretPath: viper.GetString(vaultSecretPathViperKey),
	}
}
