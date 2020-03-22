// +build unit

package chainregistry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
)

func TestApp(t *testing.T) {
	cfg := app.DefaultConfig()
	appli, err := New(cfg, nil, nil, true, nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, appli, "App should have been created")
}
