package formatters

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/zk-snarks"
)

func FormatZKSAccountResponse(account *entities.ZKSAccount) *types.ZKSAccountResponse {
	return &types.ZKSAccountResponse{
		Curve:            account.Curve,
		SigningAlgorithm: account.SigningAlgorithm,
		PublicKey:        account.PublicKey,
		Namespace:        account.Namespace,
	}
}
