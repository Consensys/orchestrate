package utils
import (
	"testing"
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