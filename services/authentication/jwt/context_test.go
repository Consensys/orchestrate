// +build unit

package jwt

import (
	"context"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	token := &jwt.Token{}
	ctx := With(context.Background(), token)
	tokenFrom := FromContext(ctx)
	assert.Equal(t, token, tokenFrom, "JWT Token should have been properly injected and extracted")
}
