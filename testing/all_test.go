//Package testing ensures the tests are gathered in a single scripts.
// This ensure we can use several time the same vault instance for
// different that would else be in differents packages.
package testing

import (
	"fmt"
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

// HashicorpInitialize returns an Initialized vault getter closure
// It ensures, it is initialized only once.
func hashicorpInitializer() *secretstore.Hashicorps {
	config := secretstore.NewConfig()
	hashicorps, err := secretstore.NewHashicorps(config)
	if err != nil {
		fmt.Printf("Error when instantiating the vault : %v", err.Error())
	}

	hashicorps.InitVault()

	return hashicorps

}

// testSecretStoreMaker returns a Test for the secretstore given a vault object
func testSecretStoreMaker(hashicorps *secretstore.Hashicorps) func(t *testing.T) {
	return func(t *testing.T) {
		var err error
		key := "secretName"
		value := "secretValue"

		err = hashicorps.Store(key, value)
		if err != nil {
			t.Errorf("Could not store the secret : %v", err.Error())
		}

		keys, err := hashicorps.List()
		if err != nil {
			t.Errorf("Could not lists the secrets : %v", err.Error())
		}
		if len(keys) != 1 || keys[0] != key {
			t.Errorf("Expected listed keys to be [%v], got %v ", key, keys)
		}

		retrievedValue, err := hashicorps.Load(key)
		if err != nil {
			t.Errorf("Could not load the secret : %v", err.Error())
		}
		if retrievedValue != value {
			t.Errorf("Expected loaded to be %v , instead got %v", value, retrievedValue)
		}

		err = hashicorps.Delete(key)
		if err != nil {
			t.Errorf("Could not delete the secret : %v", err.Error())
		}
	}
}

// testKeyStoreMaker returns a test for the keystore
func testKeyStoreMaker(hashicorps *secretstore.Hashicorps) func(t *testing.T) {
	return func(t *testing.T) {
		keystore := keystore.NewBaseKeyStore(hashicorps)

		address, err := keystore.GenerateWallet()
		if err != nil {
			t.Errorf("Error while generating a new wallet : %v", err.Error())
		}

		tx := ethtypes.NewTransaction(
			0,
			*address,
			new(big.Int).SetInt64(0),
			0,
			new(big.Int).SetInt64(0),
			[]byte{},
		)

		_, _, err = keystore.SignTx(
			(&common.Chain{}).SetID(big.NewInt(10)),
			*address,
			tx,
		)

		if err != nil {
			t.Errorf("Error while signing a transaction : %v", err.Error())
		}
	}
}

// TestAll runs all the tests as a sun tests
func TestAll(t *testing.T) {
	viper.Set("vault.token.name", "test-token")
	hashicorps := hashicorpInitializer()
	t.Run("SecretStore", testSecretStoreMaker(hashicorps))
	t.Run("KeyStore", testKeyStoreMaker(hashicorps))
}
