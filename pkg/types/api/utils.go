package api

import (
	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/types/entities"
)

func validatePrivateTxParams(protocol entities.PrivateTxManagerType, privacyGroupID string, privateFor []string) error {
	if protocol == "" {
		return errors.InvalidParameterError("field 'protocol' cannot be empty")
	}

	if privacyGroupID == "" && len(privateFor) == 0 {
		return errors.InvalidParameterError("fields 'privacyGroupId' and 'privateFor' cannot both be empty")
	}

	if len(privateFor) > 0 && privacyGroupID != "" {
		return errors.InvalidParameterError("fields 'privacyGroupId' and 'privateFor' are mutually exclusive")
	}

	return nil
}

func validateTxFromParams(from string, oneTimeKey bool) error {
	if from != "" && oneTimeKey {
		return errors.InvalidParameterError("fields 'from' and 'oneTimeKey' are mutually exclusive")
	}

	if from == "" && !oneTimeKey {
		return errors.InvalidParameterError("field 'from' is required")
	}

	return nil
}
