package main

import (
	"fmt"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/secretstore"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/keystore"
)

// @title Swagger Example API
// @version 1.0.1
func main() {

	config := secretstore.VaultConfigFromViper()
	hashicorpsSS := secretStore.NewHashicorps(config)

	awsSS = secretStore.NewAWS(7)
	tokenName := secretStore.VaultTokenFromViper()

	err = hashicorpsSS.Init(awsSS, tokenName)
	if err != nil {
		fmt.Printf(err.Error())
	}

	keystore := keystore.NewBasicKeyStore(hashicorpsSS)

}
