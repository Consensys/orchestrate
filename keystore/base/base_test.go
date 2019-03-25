package base

import (
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/hashicorp"
)

//TestSecretStore must be run along with a vault container in development mode
//It will sequentially writes a secret, list all the secrets, get the secret then delete it.
func TestKeyStore(t *testing.T) {
	config := hashicorp.NewConfig()
	hashicorpsSS, err := hashicorp.NewHashicorps(config)
	if err != nil {
		t.Errorf("Error when instantiating the vault : %v", err.Error())
	}

	err = hashicorpsSS.InitVault()
	if err != nil {
		t.Errorf("Error initializing the vault : %v", err.Error())
	}

	keystore := NewKeyStore(hashicorpsSS)

	_, err = keystore.GenerateWallet()
	if err != nil {
		t.Errorf("Error while generating a new wallet : %v", err.Error())
	}
}
