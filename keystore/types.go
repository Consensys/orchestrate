package keystore

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/hashicorp"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

// KeyStore is an interface implemented by module that are able to perform signature transactions
type KeyStore interface {
	// SignTx signs a transaction
	SignTx(chain *common.Chain, a ethcommon.Address, tx *ethtypes.Transaction) (raw []byte, txHash *ethcommon.Hash, err error)

	// SignMsg sign a message TODO: what is the EIP?
	SignMsg(a ethcommon.Address, msg string) (rsv []byte, hash *ethcommon.Hash, err error) //TODO: do not forget to add prefix

	// SignRawHash sign a bytes
	SignRawHash(a ethcommon.Address, hash []byte) (rsv []byte, err error)

	// GenerateWallet creates a wallet
	GenerateWallet() (add *ethcommon.Address, err error)
}

// NewKeyStore create new Key Store
func NewKeyStore() (KeyStore, error) {
	var s *base.KeyStore
	switch viper.GetString(secretStoreViperKey) {
	case "test":
		// Create Key Store from a Mock SecretStore
		s = base.NewKeyStore(mock.NewSecretStore())
	case "hashicorp":
		// Create an hashicorp vault object
		vault, err := hashicorp.NewHashicorps(hashicorp.NewConfig())
		if err != nil {
			return nil, err
		}

		// Create Key Store
		s = base.NewKeyStore(vault)
	default:
		// Key Store type should be one of "test", "hashicorp"
		return nil, fmt.Errorf("Invalid Store type %q", viper.GetString(secretStoreViperKey))
	}

	// Pre-Import Pkeys
	for _, pkey := range viper.GetStringSlice(secretPkeyViperKey) {
		err := s.ImportPrivateKey(pkey)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}
