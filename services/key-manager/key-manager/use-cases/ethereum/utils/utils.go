package utils

import (
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/services/key-manager/service/formatters"
	signer "github.com/ethereum/go-ethereum/signer/core"
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
