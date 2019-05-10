package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSignature(t *testing.T) {
	tests := []struct {
		sig  string
		name string
		args string
		err  error
	}{
		{"Transfer()", "Transfer", "", nil},
		{"Transfer(address[2])", "Transfer", "address[2]", nil},
		{"Transfer(uint256,address,bytes32)", "Transfer", "uint256,address,bytes32", nil},
		{"aze", "", "", errors.New("")},
	}

	for k, test := range tests {
		t.Log(k)
		assert.Equal(t, test.err == nil, IsValidSignature(test.sig))

		name, args, err := ParseSignature(test.sig)
		assert.IsType(t, test.err, err, err)
		assert.Equal(t, test.name, name)
		assert.Equal(t, test.args, args)
	}
}
