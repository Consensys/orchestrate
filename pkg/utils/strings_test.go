// +build unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortString(t *testing.T) {
	stringTest := "abcdefghijklmnopqrstuvwxyz"
	if ShortString(stringTest, 2) != "ab...yz" {
		t.Errorf("ShortHash was incorrect")
	}
	if ShortString(stringTest, 26) != stringTest {
		t.Errorf("ShortString was incorrect")
	}
}

func TestParseIArrayToStringArray(t *testing.T) {
	iArr := make([]interface{}, 2)
	iArr[0] = 0
	iArr[1] = "val1"
	sArr := ParseIArrayToStringArray(iArr)
	assert.Equal(t, "0", sArr[0])
	assert.Equal(t, "val1", sArr[1])
}
