package key

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, auth, "auth Manager should have been set")

	var a *Auth
	SetGlobalAuth(a)
	assert.Nil(t, GlobalAuth(), "Global Auth should be reset to nil")
}
