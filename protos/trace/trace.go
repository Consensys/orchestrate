package trace

import (
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

func (t *Trace) Error() string {
	return common.Errors(t.Errors).Error()
}
