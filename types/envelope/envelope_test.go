package envelope

import (
	"testing"

	"github.com/stretchr/testify/assert"
	err "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

func TestEnvelope(t *testing.T) {
	envelope := &Envelope{}
	assert.Equal(t, "", envelope.Error(), "Error message should be correct")

	envelope = &Envelope{
		Errors: []*err.Error{
			&err.Error{Code: 1, Message: "Timeout error"},
			&err.Error{Code: 0, Message: "Unknown error"},
		},
	}
	assert.Equal(t, `["Timeout error" "Unknown error"]`, envelope.Error(), "Error message should be correct")
}
