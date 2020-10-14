package formatters

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/keymanager/ethereum"
)

func FormatETHAccountResponse(account *entities.ETHAccount) *types.ETHAccountResponse {
	return &types.ETHAccountResponse{
		Address:             account.Address,
		PublicKey:           account.PublicKey,
		CompressedPublicKey: account.CompressedPublicKey,
		KeyType:             account.KeyType,
		Namespace:           account.Namespace,
	}
}
