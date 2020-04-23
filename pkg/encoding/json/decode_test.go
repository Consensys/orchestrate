// +build unit

package json

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

type unmarshalTest struct {
	in     string
	ptr    interface{}
	errMsg string
}

type mockBody struct {
	Name string   `json:"name,omitempty" validate:"required"`
	URLs []string `json:"urls,omitempty" pg:"urls,array" validate:"min=1,unique,dive,url"`
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

func TestUnmarshalBody(t *testing.T) {
	testSuite := []struct {
		name           string
		body           func() []byte
		input          interface{}
		expectedOutput interface{}
		expectedError  error
	}{
		{
			name:           "unknown field",
			body:           func() []byte { return []byte(`{"unknownField":"error"}`) },
			input:          &mockBody{},
			expectedOutput: &mockBody{},
			expectedError:  errors.InvalidFormatError("json: unknown field \"unknownField\"").ExtendComponent(component),
		},
		{
			name: "twice same URL field",
			body: func() []byte {
				body, _ := Marshal(&mockBody{
					Name: "testName",
					URLs: []string{"http://test.com", "http://test.com"},
				})
				return body
			},
			input:          &mockBody{},
			expectedOutput: &mockBody{Name: "testName", URLs: []string{"http://test.com", "http://test.com"}},
			expectedError:  errors.InvalidParameterError("invalid body, with: field validation for 'URLs' failed on the 'unique' tag").ExtendComponent(component),
		},
	}

	for _, test := range testSuite {
		test := test // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(test.name, func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other

			err := UnmarshalBody(bytes.NewReader(test.body()), test.input)

			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.input, test.expectedOutput)
		})
	}
}
