// +build unit

package jwt

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	token := "myToken"
	ctx := With(context.Background(), token)
	tokenFrom := FromContext(ctx)
	assert.Equal(t, token, tokenFrom)
}
