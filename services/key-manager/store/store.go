package store

import (
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	types "github.com/ConsenSys/orchestrate/pkg/types/keymanager/ethereum"
)

//go:generate mockgen -source=store.go -destination=mocks/mock.go -package=mocks

type Vault interface {
	ETHCreateAccount(namespace string) (*entities.ETHAccount, error)
	ETHImportAccount(namespace, privateKey string) (*entities.ETHAccount, error)
	ETHSign(address string, namespace, data string) (string, error)
	ETHSignTransaction(address string, request *types.SignETHTransactionRequest) (string, error)
	ETHSignQuorumPrivateTransaction(address string, request *types.SignQuorumPrivateTransactionRequest) (string, error)
	ETHSignEEATransaction(address string, request *types.SignEEATransactionRequest) (string, error)
	ETHListAccounts(namespace string) ([]string, error)
	ETHListNamespaces() ([]string, error)
	ETHGetAccount(address string, namespace string) (*entities.ETHAccount, error)

	ZKSCreateAccount(namespace string) (*entities.ZKSAccount, error)
	ZKSListNamespaces() ([]string, error)
	ZKSSign(pubKey string, namespace, data string) (string, error)
	ZKSListAccounts(namespace string) ([]string, error)
	ZKSGetAccount(pubKey string, namespace string) (*entities.ZKSAccount, error)

	HealthCheck() error
}
