package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

func TestEnvelope(t *testing.T) {
	envelope := &Envelope{}
	assert.Equal(t, "", envelope.Error(), "Error message should be correct")

	envelope = &Envelope{
		Errors: []*common.Error{
			&common.Error{Type: 1, Message: "Timeout error"},
			&common.Error{Type: 0, Message: "Unknown error"},
		},
	}
	assert.Equal(t, `2 error(s): ["Error #1: Timeout error" "Error #0: Unknown error"]`, envelope.Error(), "Error message should be correct")
}
