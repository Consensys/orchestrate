package envelope

import (
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

func (t *Envelope) Error() string {
	return common.Errors(t.Errors).Error()
}
