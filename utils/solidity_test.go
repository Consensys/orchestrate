package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	err "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/error"
)

func TestGetSignature(t *testing.T) {
	tests := []struct {
		sig  string
		name string
		args string
		err  bool
	}{
		{"Transfer()", "Transfer", "", false},
		{"Transfer(address[2])", "Transfer", "address[2]", false},
		{"Transfer(uint256,address,bytes32)", "Transfer", "uint256,address,bytes32", false},
		{"aze", "", "", true},
	}

	for k, test := range tests {
		t.Log(k)
		assert.Equal(t, !test.err, IsValidSignature(test.sig))

		name, args, e := ParseSignature(test.sig)
		assert.Equal(t, test.name, name)
		assert.Equal(t, test.args, args)
		assert.Equal(t, test.err, e != nil, e)
		if e != nil {
			ie, ok := e.(*err.Error)
			assert.True(t, ok, "Error should cast to internal error")
			assert.Equal(t, "utils", ie.GetComponent(), "Error component should be valid")
		}
	}
}
