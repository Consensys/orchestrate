// +build unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Field struct{}

type Test struct {
	Field1 *Field
	Field2 *Field
}

func TestExtractField(t *testing.T) {
	testField := &Field{}
	srv := &Test{
		Field1: testField,
	}

	field, err := ExtractField(srv)
	require.NoError(t, err)
	assert.Equal(t, testField, field)

	// Empty
	srv = &Test{}
	_, err = ExtractField(srv)
	require.Error(t, err)

	// Empty
	srv = &Test{
		Field1: &Field{},
		Field2: &Field{},
	}
	_, err = ExtractField(srv)
	require.Error(t, err)
}
