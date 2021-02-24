package formatters

import (
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	types "github.com/ConsenSys/orchestrate/pkg/types/keymanager/zk-snarks"
)

func FormatZKSAccountResponse(account *entities.ZKSAccount) *types.ZKSAccountResponse {
	return &types.ZKSAccountResponse{
		Curve:            account.Curve,
		SigningAlgorithm: account.SigningAlgorithm,
		PublicKey:        account.PublicKey,
		Namespace:        account.Namespace,
	}
}
