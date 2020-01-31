package utils

import (
	"context"
	"testing"
)

func TestAuthorizationFromContext(t *testing.T) {
	test := AuthorizationFromContext(context.Background())
	t.Log(test)
}

func TestApiKeyFromContext(t *testing.T) {
	test := APIKeyFromContext(context.Background())
	t.Log(test)
}
