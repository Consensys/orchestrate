package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	b, err := Marshal(true)
	assert.Nil(t, err, "Marshal should not error on true")
	assert.Equal(t, []byte{0x74, 0x72, 0x75, 0x65}, b, "Marshal true should be correct")
}

func TestMarshalError(t *testing.T) {
	b, err := Marshal(make(chan int))
	assert.Nil(t, b, "Marshal should be nil")
	assert.Error(t, err, "Should get error")
}
