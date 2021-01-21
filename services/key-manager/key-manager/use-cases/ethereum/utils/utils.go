package utils

import (
	"fmt"

	signer "github.com/ethereum/go-ethereum/signer/core"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/service/formatters"
)

func GetEIP712EncodedData(typedData *signer.TypedData) (string, error) {
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return "", errors.InvalidParameterError("invalid typed data message").AppendReason(err.Error())
	}

	domainSeparatorHash, err := typedData.HashStruct(formatters.DomainLabel, typedData.Domain.Map())
	if err != nil {
		return "", errors.InvalidParameterError("invalid domain separator").AppendReason(err.Error())
	}

	return fmt.Sprintf("\x19\x01%s%s", domainSeparatorHash, typedDataHash), nil
}
