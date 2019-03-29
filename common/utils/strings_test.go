package utils
import (
	"testing"
)

func TestShortHash(t *testing.T) {
	stringTest := "abcdefghijklmnopqrstuvwxyz"
	if ShortHash(stringTest, 2) != "ab...yz" {
		t.Errorf("ShortHash was incorrect")
	}
	if ShortHash(stringTest, 26) != stringTest {
		t.Errorf("ShortHash was incorrect")
	}
}