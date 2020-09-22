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
	iArr := make([]interface{}, 6)
	iArr[0] = 0
	iArr[1] = []int{0, 1}

	iArr[2] = "0x123"
	iArr[3] = []string{"0x123", "0x213"}

	iArr[4] = false
	iArr[5] = []bool{false, true}
	sArr, _ := ParseIArrayToStringArray(iArr)

	assert.Equal(t, "0", sArr[0])
	assert.Equal(t, `["0","1"]`, sArr[1])

	assert.Equal(t, "0x123", sArr[2])
	assert.Equal(t, `["0x123","0x213"]`, sArr[3])

	assert.Equal(t, "false", sArr[4])
	assert.Equal(t, `["false","true"]`, sArr[5])
}
