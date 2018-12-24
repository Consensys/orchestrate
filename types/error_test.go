package types

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := Error{fmt.Errorf("test"), 127}

	if err.Error() != "test" {
		t.Errorf("Error: expected %v but got %v", "test", err.Error())
	}
}
