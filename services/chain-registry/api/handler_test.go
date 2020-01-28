package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

func TestUnmarshalBody(t *testing.T) {
	testSuite := []struct {
		name           string
		body           func() []byte
		input          interface{}
		expectedOutput interface{}
		expectecError  error
	}{
		{
			name:           "unknown field",
			body:           func() []byte { return []byte(`{"unknownField":"error"}`) },
			input:          &PostRequest{},
			expectedOutput: &PostRequest{},
			expectecError:  errors.FromError(fmt.Errorf("json: unknown field \"unknownField\"")).ExtendComponent(component),
		},
		{
			name: "twice same URL field",
			body: func() []byte {
				body, _ := json.Marshal(&PostRequest{
					Name: "testName",
					URLs: []string{"http://test.com", "http://test.com"},
				})
				return body
			},
			input:          &PostRequest{},
			expectedOutput: &PostRequest{Name: "testName", URLs: []string{"http://test.com", "http://test.com"}},
			expectecError:  errors.FromError(fmt.Errorf("invalid body")).ExtendComponent(component),
		},
	}

	for _, test := range testSuite {
		test := test // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(test.name, func(t *testing.T) {
			t.Parallel() // marks each test case as capable of running in parallel with each other

			err := UnmarshalBody(bytes.NewReader(test.body()), test.input)

			assert.Equal(t, test.expectecError, err, "should get same error")
			assert.Equal(t, test.input, test.expectedOutput, "should unmarshal without error")
		})
	}

}
