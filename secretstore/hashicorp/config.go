package hashicorp

import (
	"fmt"

	vault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/aws"
)

func init() {
	viper.SetDefault(vaultURIViperKey, vaultURIDefault)
	viper.BindEnv(vaultURIViperKey, vaultURIEnv)
	viper.SetDefault(vaultTokenNameViperKey, vaultTokenNameEnv)
	viper.BindEnv(vaultTokenNameViperKey, vaultTokenNameEnv)
	viper.SetDefault(vaultTokenViperKey, vaultTokenEnv)
	viper.BindEnv(vaultTokenViperKey, vaultTokenEnv)
	viper.SetDefault(vaultUnsealKeyViperKey, vaultUnsealKeyEnv)
	viper.BindEnv(vaultUnsealKeyViperKey, vaultUnsealKeyEnv)
}

var (
	vaultURIFlag     = "vault-uri"
	vaultURIViperKey = "vault.uri"
	vaultURIDefault  = "http://127.0.0.1:8200"
	vaultURIEnv      = "VAULT_URI"

	vaultTokenNameFlag     = "vault-token-name"
	vaultTokenNameViperKey = "vault.token.name"
	vaultTokenNameDefault  = ""
	vaultTokenNameEnv      = "VAULT_TOKEN_NAME"

	vaultTokenFlag     = "vault-token"
	vaultTokenViperKey = "vault.token"
	vaultTokenDefault  = ""
	vaultTokenEnv      = "VAULT_TOKEN"

	vaultUnsealKeyFlag     = "vault-unseal-key"
	vaultUnsealKeyViperKey = "vault.unseal.key"
	vaultUnsealKeyDefault  = ""
	vaultUnsealKeyEnv      = "VAULT_UNSEAL_KEY"
)

// InitFlags register flags for hashicorp vault
func InitFlags(f *pflag.FlagSet) {
	VaultURI(f)
	VaultToken(f)
	VaultUnsealKey(f)
	VaultTokenName(f)
}

// VaultURI register a flag for vault server address
func VaultURI(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret vault URI 
Environment variable: %q`, vaultURIEnv)
	f.String(vaultURIFlag, vaultURIDefault, desc)
	viper.BindPFlag(vaultURIViperKey, f.Lookup(vaultURIFlag))
}

// VaultToken register a flag for vault server address
func VaultToken(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp secret vault authentication token
Environment variable: %q`, vaultTokenEnv)
	f.String(vaultTokenFlag, vaultTokenDefault, desc)
	viper.BindPFlag(vaultTokenViperKey, f.Lookup(vaultTokenFlag))
}

// VaultUnsealKey registers a flag for the value of the vault unseal key
func VaultUnsealKey(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp vault unsealing key
Environment variable: %q`, vaultTokenEnv)
	f.String(vaultUnsealKeyFlag, vaultUnsealKeyDefault, desc)
	viper.BindPFlag(vaultUnsealKeyViperKey, f.Lookup(vaultUnsealKeyFlag))
}

// VaultTokenName register a flag for vault server address
func VaultTokenName(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Hashicorp token name on AWS
Environment variable: %q`, vaultTokenNameEnv)
	f.String(vaultTokenNameFlag, vaultTokenNameDefault, desc)
	viper.BindPFlag(vaultTokenNameViperKey, f.Lookup(vaultURIFlag))
}

// NewConfig icreates vault configuration from viper
func NewConfig() *vault.Config {
	config := vault.DefaultConfig()
	config.Address = viper.GetString(vaultURIViperKey)
	return config
}

// AutoInit will try to Init the vault directly or FetchFromAws
func AutoInit(hashicorps *Hashicorps) (err error) {
	tokenName := viper.GetString(vaultTokenNameViperKey)
	log.WithFields(log.Fields{
		"aws.secret": tokenName,
	}).Debugf("hashicorp: auto-initiliazing vault (credentials from AWS)")

	// Create AWS object to retrieve Vault credentials
	awsSS := aws.NewAWS(7)

	// Initialize Vault
	err = hashicorps.InitVault()
	if err != nil {
		log.WithError(err).Debugf("hashicorp: failed to init")

		// Probably Vault is already unsealed so we retrieve credentials from AWS
		secret, err := awsSS.Load(tokenName)
		if err != nil {
			log.WithError(err).Errorf("hashicorp: failed to load credentials from AWS")
			return err
		}

		err = hashicorps.creds.fromEncoded(secret)
		if err != nil {
			log.WithError(err).Errorf("hashicorp: failed to decode credentials")
			return err
		}

		log.WithFields(log.Fields{
			"vault.token": hashicorps.creds.Token,
		}).Warnf("hashicorp: !!!WARNING this message is meant for debugging purpose and should be removed ASAP!!!")
		hashicorps.SetToken(hashicorps.creds.Token)
	} else {
		// Vault has been properly unsealed so we push credentials on AWS
		encoded, err := hashicorps.creds.encode()
		if err != nil {
			return err
		}

		err = awsSS.Store(tokenName, encoded)
		if err != nil {
			log.WithError(err).Errorf("hashicorp: failed to store credentials in AWS")
			return fmt.Errorf("Could not send credentials to AWS: %v", err)
		}
	}
	return nil
}
