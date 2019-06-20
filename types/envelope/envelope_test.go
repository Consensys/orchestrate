package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
)

func TestEnvelope(t *testing.T) {
	envelope := &Envelope{}
	assert.Equal(t, "", envelope.Error(), "Error message should be correct")

	envelope = &Envelope{
		Errors: []*common.Error{
			&common.Error{Code: 1, Message: "Timeout error"},
			&common.Error{Code: 0, Message: "Unknown error"},
		},
	}
	assert.Equal(t, `2 error(s): ["Error #1: Timeout error" "Error #0: Unknown error"]`, envelope.Error(), "Error message should be correct")
}
