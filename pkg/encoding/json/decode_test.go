// +build unit

package json

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/consensys/orchestrate/pkg/errors"
)

type mockBody struct {
	Name string   `json:"name,omitempty" validate:"required"`
	URLs []string `json:"urls,omitempty" pg:"urls,array" validate:"min=1,unique,dive,url"`
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
			expectedError:  errors.InvalidFormatError("failed to decode request body (json: unknown field \"unknownField\")"),
		},
		{
			name: "twice same URL field",
			body: func() []byte {
				body, _ := json.Marshal(&mockBody{
					Name: "testName",
					URLs: []string{"http://test.com", "http://test.com"},
				})
				return body
			},
			input:          &mockBody{},
			expectedOutput: &mockBody{Name: "testName", URLs: []string{"http://test.com", "http://test.com"}},
			expectedError:  errors.InvalidParameterError("invalid body (field validation for 'URLs' failed on the 'unique' tag)"),
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
