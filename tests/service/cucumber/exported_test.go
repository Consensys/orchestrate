package cucumber

import (
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// Init(context.Background())
	// assert.NotNil(t, GlobalOptions(), "Global should have been set") no-lint:gocritic

	var o *godog.Options
	SetGlobalOptions(o)
	assert.Nil(t, GlobalOptions(), "Global should be reset to nil")
}
