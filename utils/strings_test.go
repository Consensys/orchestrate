package utils
import (
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git"
	"testing"
)

func TestShortString(t *testing.T) {
	stringTest := "abcdefghijklmnopqrstuvwxyz"
	if pkg.ShortString(stringTest, 2) != "ab...yz" {
		t.Errorf("ShortHash was incorrect")
	}
	if pkg.ShortString(stringTest, 26) != stringTest {
		t.Errorf("ShortString was incorrect")
	}
}