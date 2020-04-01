// +build unit

package handlers

import (
	"context"
	"testing"
)

func TestInit(t *testing.T) {
	Init(context.Background())
}
