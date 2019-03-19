package keystore

import (
	"github.com/spf13/cobra"
	"testing"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore"
)

//TestSecretStore must be run along with a vault container in development mode
//It will sequentially writes a secret, list all the secrets, get the secret then delete it.
func TestKeyStore(t *testing.T) {

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   func(cmd *cobra.Command, args []string){},
	}

	secretstore.VaultURI(runCmd.Flags())

	config := secretstore.VaultConfigFromViper()
	hashicorpsSS, err := secretstore.NewHashicorps(config)
	if err != nil {
		t.Errorf("Error when instantiating the vault : %v", err.Error())
	}

	err = hashicorpsSS.InitVault()
	if err != nil {
		t.Errorf("Error initializing the vault : %v", err.Error())
	}

	keystore := NewBaseKeyStore(hashicorpsSS)

	_, err = keystore.GenerateWallet()
	if err != nil {
		t.Errorf("Error while generating a new wallet : %v", err.Error())
	}

}