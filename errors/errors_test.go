package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigError(t *testing.T) {
	assert.Equal(t, []byte{0x10, 0x00}, ConfigError(fmt.Errorf("test")).GetCode(), "ConfigError code should be correct")
}

func TestEncodingError(t *testing.T) {
	assert.Equal(t, []byte{0x20, 0x00}, EncodingError(fmt.Errorf("test")).GetCode(), "EncodingError code should be correct")
}

func TestDataError(t *testing.T) {
	assert.Equal(t, []byte{0x30, 0x00}, DataError(fmt.Errorf("test")).GetCode(), "DataError code should be correct")
}

func TestSolidityError(t *testing.T) {
	assert.Equal(t, []byte{0x31, 0x00}, SolidityError(fmt.Errorf("test")).GetCode(), "SolidityError code should be correct")
}

func TestInvalidSigError(t *testing.T) {
	assert.Equal(t, []byte{0x31, 0x01}, InvalidSigError("test").GetCode(), "InvalidSigError code should be correct")
}
