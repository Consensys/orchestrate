package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type unmarshalTest struct {
	in     string
	ptr    interface{}
	errMsg string
}

var unmarshalTests = []unmarshalTest{
	{in: `true`, ptr: new(bool)},
	{in: `1`, ptr: new(int)},
	{in: `{"x": 1}`, ptr: new(bool), errMsg: "json: cannot unmarshal object into Go value of type bool"},
}

func TestUnmarshal(t *testing.T) {
	for _, test := range unmarshalTests {
		in := []byte(test.in)
		err := errors.FromError(Unmarshal(in, test.ptr))
		if test.errMsg == "" {
			assert.Nil(t, err, "Unmarshal should not error")
		} else {
			assert.Error(t, err, "Unmarshal should error")
			assert.Equal(t, "encoding.json", err.GetComponent(), "Error code should be correct")
			assert.Equal(t, test.errMsg, err.GetMessage(), "Error message should be correct")
		}
	}
}
