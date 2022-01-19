package utils

import (
	"testing"
)

func TestToPtr_SimpleTypes(t *testing.T) {
	_ = ToPtr("a").(*string)
	_ = ToPtr(1).(*int)
	_ = ToPtr([]string{"a"}).(*[]string)
}
