package infra

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore"
)

var (
	vaultAccountFlag     = "vault-account"
	vaultAccountViperKey = "vault.accounts"
	vaultAccountDefault  = []string{
		"56202652FDFFD802B7252A456DBD8F3ECC0352BBDE76C23B40AFE8AEBD714E2E",
		"5FBB50BFF6DFAD35C4A374C9237BA2F7EAED9C6868E0108CB259B62D68029B1A",
		"86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC",
		"DD614C3B343E1B6DBD1B2811D4F146CC90337DEEF96AB97C353578E871B19D5E",
		"425D92F63A836F890F1690B34B6A25C2971EF8D035CD8EA8592FD1069BD151C6",
		"C4B172E72033581BC41C36FA0448FCF031E9A31C4A3E300E541802DFB7248307",
		"706CC0876DA4D52B6DCE6F5A0FF210AEFCD51DE9F9CFE7D1BF7B385C82A06B8C",
		"1476C66DE79A57E8AB4CADCECCBE858C99E5EDF3BFFEA5404B15322B5421E18C",
		"A2426FE76ECA2AA7852B95A2CE9CC5CC2BC6C05BB98FDA267F2849A7130CF50D",
		"41B9C5E497CFE6A1C641EFCA314FF84D22036D1480AF5EC54558A5EDD2FEAC03",
	}
	vaultAccountEnv = "VAULT_ACCOUNT"
)

// VaultAccounts register flag for Vault accounts
func VaultAccounts(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address of Kafka server to connect to.
Environment variable: %q`, vaultAccountEnv)
	f.StringSlice(vaultAccountFlag, vaultAccountDefault, desc)
	viper.BindPFlag(vaultAccountViperKey, f.Lookup(vaultAccountFlag))
	viper.BindEnv(vaultAccountViperKey, vaultAccountEnv)
}

func initSigner(infra *Infra) error {
	// h.cfg.Vault.Accounts

	config := secretstore.VaultConfigFromViper()
	hashi, err := secretstore.NewHashicorps(config)
	if err != nil {
		return err
	}

	if !ManualInit(hashi) {
		err := AutoInit(hashi)
		if err != nil {
			return err
		}
	}

	infra.SecretStore = hashi
	infra.KeyStore = keystore.NewBaseKeyStore(hashi)

	log.Infof("infra-signer: ready")
	return nil
}

// ManualInit (UNSAFE) will Init the vault with the value directly passed 
// Ignore if the token variable is not set
func ManualInit(hash *secretstore.Hashicorps) bool {

	token := secretstore.VaultTokenFromViper()
	if token != "" {

		hash.Client.SetToken(token)
		unsealKey := secretstore.VaultUnsealKeyFromViper()
		 
		if unsealKey != "" {
			hash.Unseal(unsealKey)
		}

		return true
	}

	return false
}

// AutoInit will try to Init the vault directly or FetchFromAws
func AutoInit(hash *secretstore.Hashicorps) (error) {

	tokenName := secretstore.VaultTokenFromViper()
	awsSS := secretstore.NewAWS(7)
	
	err := hash.InitVault(awsSS, tokenName)
	if err != nil {
		err2 := hash.InitFromAWS(awsSS, tokenName)
		if err2 != nil {
			return fmt.Errorf("Got %v when trying to init the vault then when trying to fetch on AWS got %v", err.Error(), err2.Error())
		}
	}

	return nil
}

// WriteInitializedKey will add the keys in the vault
func WriteInitializedKey(bks *keystore.BaseKeyStore) (err error) {

	accounts := viper.GetStringSlice("vault.accounts")
	for _, account := range accounts {
		err = bks.ImportPrivateKey(account)
		if err != nil {
			return err
		}
	}

	return nil

}
