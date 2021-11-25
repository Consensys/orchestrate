package api

import (
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func validatePrivateTxParams(protocol entities.PrivateTxManagerType, privateFrom, privacyGroupID string, privateFor []string) error {
	if protocol == "" {
		return errors.InvalidParameterError("field 'protocol' cannot be empty")
	}

	if protocol != entities.TesseraChainType && privateFrom == "" {
		return errors.InvalidParameterError("fields 'privateFrom' cannot be empty")
	}

	if privacyGroupID == "" && len(privateFor) == 0 {
		return errors.InvalidParameterError("fields 'privacyGroupId' and 'privateFor' cannot both be empty")
	}

	if len(privateFor) > 0 && privacyGroupID != "" {
		return errors.InvalidParameterError("fields 'privacyGroupId' and 'privateFor' are mutually exclusive")
	}

	return nil
}

func validateTxFromParams(from *ethcommon.Address, oneTimeKey bool) error {
	if from != nil && oneTimeKey {
		return errors.InvalidParameterError("fields 'from' and 'oneTimeKey' are mutually exclusive")
	}

	if from == nil && !oneTimeKey {
		return errors.InvalidParameterError("field 'from' is required")
	}

	return nil
}
