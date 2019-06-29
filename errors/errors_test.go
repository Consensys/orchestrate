package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWarning(t *testing.T) {
	assert.Equal(t, uint64(4096), Warning("test").GetCode(), "Warning code should be correct")
}

func TestConnectionError(t *testing.T) {
	assert.Equal(t, uint64(32768), ConnectionError("test").GetCode(), "ConnectionError code should be correct")
}

func TestConfigError(t *testing.T) {
	assert.Equal(t, uint64(983040), ConfigError("test").GetCode(), "ConfigError code should be correct")
}

func TestDataError(t *testing.T) {
	assert.Equal(t, uint64(270336), DataError("test").GetCode(), "DataError code should be correct")
}

func TestEncodingError(t *testing.T) {
	e := EncodingError("test")
	assert.Equal(t, uint64(270592), e.GetCode(), "EncodingError code should be correct")
	assert.True(t, IsDataError(e), "EncodingError should be a data error")
}

func TestSolidityError(t *testing.T) {
	e := SolidityError("test")
	assert.Equal(t, uint64(270848), e.GetCode(), "SolidityError code should be correct")
	assert.True(t, IsDataError(e), "SolidityError should be a data error")
}

func TestInvalidSigError(t *testing.T) {
	e := InvalidSigError("test")
	assert.Equal(t, uint64(270849), e.GetCode(), "InvalidSigError code should be correct")
	assert.True(t, IsDataError(e), "InvalidSigError should be a data error")
	assert.True(t, IsSolidityError(e), "IsSolidityError should be a data error")
}

func TestInvalidFormatError(t *testing.T) {
	e := InvalidFormatError("test")
	assert.Equal(t, uint64(271104), e.GetCode(), "InvalidFormatError code should be correct")
	assert.True(t, IsDataError(e), "InvalidFormatError should be a data error")
}

func TestInvalidFormatErrorf(t *testing.T) {
	assert.Equal(t, uint64(271104), InvalidFormatErrorf("test %q", "value").GetCode(), "InvalidFormatErrorf code should be correct")
}

func TestIs(t *testing.T) {
	assert.True(t, is(271120, dataErrCode), "Hex 42310 should be a data error")
	assert.False(t, is(dataErrCode, solidityErrCode), "Data error should not be a solidity error")
	assert.False(t, is(dataErrCode, 0), "Hex 00000 should not be a data error")
	assert.False(t, is(0, dataErrCode), "Hex 00000 should not be a data error")
	assert.False(t, is(275216, dataErrCode), "Hex 43310 should not be a data error")
}
