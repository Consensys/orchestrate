package store

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
)

//go:generate mockgen -source=store.go -destination=mocks/mock.go -package=mocks

type Vault interface {
	ETHCreateAccount(namespace string) (*entities.ETHAccount, error)
	ETHImportAccount(namespace, privateKey string) (*entities.ETHAccount, error)
	ETHSign(address string, namespace, data string) (string, error)
	ETHSignTransaction(address string, request *types.SignETHTransactionRequest) (string, error)
	ETHSignQuorumPrivateTransaction(address string, request *types.SignQuorumPrivateTransactionRequest) (string, error)
	ETHSignEEATransaction(address string, request *types.SignEEATransactionRequest) (string, error)
	HealthCheck() error
}
