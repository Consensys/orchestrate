package utils

import (
	"context"
	"testing"
)

func TestAuthorizationFromContext(t *testing.T) {
	test := AuthorizationFromContext(context.Background())
	t.Log(test)
}
