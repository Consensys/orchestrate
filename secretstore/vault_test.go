package secretstore

import (
	"github.com/spf13/cobra"
	"testing"
)

//TestSecretStore must be run along with a vault container in development mode
//It will sequentially writes a secret, list all the secrets, get the secret then delete it.
func TestSecretStore(t *testing.T) {

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run application",
		Run:   func(cmd *cobra.Command, args []string){},
	}

	VaultURI(runCmd.Flags())

	config := VaultConfigFromViper()
	hashicorpsSS, err := NewHashicorps(config)
	if err != nil {
		t.Errorf("Error when instantiating the vault : %v", err.Error())
	}

	hashicorpsSS.creds.FetchFromVaultInit(hashicorpsSS.Client)

	err = hashicorpsSS.creds.Unseal(hashicorpsSS.Client)
	if err != nil {
		t.Errorf("Error unsealing vault : %v (the UnsealKey sent was %v)", err.Error(), hashicorpsSS.creds.Keys)
	}
	hashicorpsSS.creds.AttachTo(hashicorpsSS.Client)

	key := "secretName"
	value := "secretValue"

	err = hashicorpsSS.Store(key, value)
	if err != nil {
		t.Errorf("Could not store the secret : %v", err.Error())
	}

	keys, err := hashicorpsSS.List()
	if err != nil {
		t.Errorf("Could not lists the secrets : %v", err.Error())
	}
	if len(keys) != 1 || keys[0] != key {
		t.Errorf("Expected listed keys to be [%v], got %v ", key, keys)
	}

	retrievedValue, err := hashicorpsSS.Load(key)
	if err != nil {
		t.Errorf("Could not load the secret : %v", err.Error())
	}
	if retrievedValue != value {
		t.Errorf("Expected loaded to be %v , instead got %v", value, retrievedValue)
	}

	err = hashicorpsSS.Delete(key)
	if err != nil {
		t.Errorf("Could not delete the secret : %v", err.Error())
	}

}