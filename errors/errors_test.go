package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigError(t *testing.T) {
	assert.Equal(t, []byte{0x10, 0x00}, ConfigError("test").GetCode(), "ConfigError code should be correct")
}

func TestEncodingError(t *testing.T) {
	assert.Equal(t, []byte{0x20, 0x00}, EncodingError("test").GetCode(), "EncodingError code should be correct")
}

func TestDataError(t *testing.T) {
	assert.Equal(t, []byte{0x30, 0x00}, DataError("test").GetCode(), "DataError code should be correct")
}

func TestSolidityError(t *testing.T) {
	assert.Equal(t, []byte{0x31, 0x00}, SolidityError("test").GetCode(), "SolidityError code should be correct")
}

func TestInvalidSigError(t *testing.T) {
	assert.Equal(t, []byte{0x31, 0x01}, InvalidSigError("test").GetCode(), "InvalidSigError code should be correct")
}

func TestInvalidFormatError(t *testing.T) {
	assert.Equal(t, []byte{0x32, 0x00}, InvalidFormatError("test").GetCode(), "InvalidFormatError code should be correct")
}

func TestInvalidFormatErrorf(t *testing.T) {
	assert.Equal(t, []byte{0x32, 0x00}, InvalidFormatErrorf("test %q", "value").GetCode(), "InvalidFormatErrorf code should be correct")
}
