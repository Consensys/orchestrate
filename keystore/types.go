package keystore

import (
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/keystore/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/hashicorp"
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
	switch viper.GetString(secretStoreViperKey) {
	case "test":
		s := mock.NewKeyStore()
		for _, pkey := range viper.GetStringSlice(secretPkeyViperKey) {
			err := s.ImportPrivateKey(pkey)
			if err != nil {
				return nil, err
			}
		}
		return s, nil
	case "hashicorp":
		vault, err := hashicorp.NewHashicorps(hashicorp.NewConfig())
		if err != nil {
			return nil, err
		}

		// Initialize hashicorp Vault
		err = hashicorp.AutoInit(vault)
		if err != nil {
			return nil, err
		}

		s := base.NewKeyStore(vault)

		for _, pkey := range viper.GetStringSlice(secretPkeyViperKey) {
			err := s.ImportPrivateKey(pkey)
			if err != nil {
				return nil, err
			}
		}

		return s, nil
	default:
		return nil, fmt.Errorf("Invalid Key Store %q", viper.GetString(secretStoreViperKey))
	}
}
