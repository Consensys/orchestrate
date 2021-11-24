package formatters

import (
	"net/http"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"

	types "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
)

func FormatRegisterFaucetRequest(request *types.RegisterFaucetRequest) *entities.Faucet {
	return &entities.Faucet{
		Name:            request.Name,
		ChainRule:       request.ChainRule,
		CreditorAccount: ethcommon.HexToAddress(request.CreditorAccount).Hex(),
		MaxBalance:      request.MaxBalance,
		Amount:          request.Amount,
		Cooldown:        request.Cooldown,
	}
}

func FormatUpdateFaucetRequest(request *types.UpdateFaucetRequest, uuid string) *entities.Faucet {
	return &entities.Faucet{
		UUID:            uuid,
		Name:            request.Name,
		ChainRule:       request.ChainRule,
		CreditorAccount: ethcommon.HexToAddress(request.CreditorAccount).Hex(),
		MaxBalance:      request.MaxBalance,
		Amount:          request.Amount,
		Cooldown:        request.Cooldown,
	}
}

func FormatFaucetResponse(faucet *entities.Faucet) *types.FaucetResponse {
	return &types.FaucetResponse{
		UUID:            faucet.UUID,
		Name:            faucet.Name,
		TenantID:        faucet.TenantID,
		ChainRule:       faucet.ChainRule,
		CreditorAccount: faucet.CreditorAccount,
		MaxBalance:      faucet.MaxBalance,
		Amount:          faucet.Amount,
		Cooldown:        faucet.Cooldown,
		CreatedAt:       faucet.CreatedAt,
		UpdatedAt:       faucet.UpdatedAt,
	}
}

func FormatFaucetFilters(req *http.Request) (*entities.FaucetFilters, error) {
	filters := &entities.FaucetFilters{}

	qNames := req.URL.Query().Get("names")
	if qNames != "" {
		filters.Names = strings.Split(qNames, ",")
	}

	qChainRule := req.URL.Query().Get("chain_rule")
	if qChainRule != "" {
		filters.ChainRule = qChainRule
	}

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, err
	}

	return filters, nil
}
