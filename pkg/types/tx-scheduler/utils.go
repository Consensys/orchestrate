package txschedulertypes

import "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

func validatePrivateTxParams(protocol, privacyGroupID string, privateFor []string) error {
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
