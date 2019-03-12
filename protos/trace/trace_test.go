package ethereum

import (
	"testing"

	"github.com/stretchr/testify/assert"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

func TestTrace(t *testing.T) {
	trace := &Trace{}
	assert.Equal(t, "", trace.Error(), "Error message should be correct")

	trace = &Trace{
		Errors: []*common.Error{
			&common.Error{Type: 1, Message: "Timeout error"},
			&common.Error{Type: 0, Message: "Unknown error"},
		},
	}
	assert.Equal(t, `2 error(s): ["Error #1: Timeout error" "Error #0: Unknown error"]`, trace.Error(), "Error message should be correct")
}
