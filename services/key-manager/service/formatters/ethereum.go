package formatters

import (
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	types "github.com/ConsenSys/orchestrate/pkg/types/keymanager/ethereum"
	"github.com/ethereum/go-ethereum/common/math"
	signer "github.com/ethereum/go-ethereum/signer/core"
)

const DomainLabel = "EIP712Domain"

func FormatETHAccountResponse(account *entities.ETHAccount) *types.ETHAccountResponse {
	return &types.ETHAccountResponse{
		Address:             account.Address,
		PublicKey:           account.PublicKey,
		CompressedPublicKey: account.CompressedPublicKey,
		Namespace:           account.Namespace,
	}
}

func FormatSignTypedDataRequest(request *types.SignTypedDataRequest) *signer.TypedData {
	typedData := &signer.TypedData{
		Types: signer.Types{
			DomainLabel: []signer.Type{
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "version", Type: "string"},
			},
		},
		PrimaryType: request.MessageType,
		Domain: signer.TypedDataDomain{
			Name:              request.DomainSeparator.Name,
			Version:           request.DomainSeparator.Version,
			ChainId:           math.NewHexOrDecimal256(request.DomainSeparator.ChainID),
			VerifyingContract: request.DomainSeparator.VerifyingContract,
			Salt:              request.DomainSeparator.Salt,
		},
		Message: request.Message,
	}

	for typeName, requestTypes := range request.Types {
		var typesDefinition []signer.Type
		for _, typeDefRequest := range requestTypes {
			typesDefinition = append(typesDefinition, signer.Type{
				Name: typeDefRequest.Name,
				Type: typeDefRequest.Type,
			})
		}
		typedData.Types[typeName] = typesDefinition
	}

	if request.DomainSeparator.VerifyingContract != "" {
		typedData.Types[DomainLabel] = append(typedData.Types[DomainLabel], signer.Type{Name: "verifyingContract", Type: "address"})
	}

	if request.DomainSeparator.Salt != "" {
		typedData.Types[DomainLabel] = append(typedData.Types[DomainLabel], signer.Type{Name: "salt", Type: "string"})
	}

	return typedData
}
