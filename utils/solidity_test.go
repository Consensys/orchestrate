package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

		name, args, err := ParseSignature(test.sig)
		assert.Equal(t, test.err, err != nil, err)
		assert.Equal(t, test.name, name)
		assert.Equal(t, test.args, args)
	}
}
