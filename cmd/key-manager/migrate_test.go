// +build unit

package keymanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrate(t *testing.T) {
	migrateCmd := newMigrateCmd()
	assert.NotNil(t, migrateCmd, "run cmd should not be nil")
}
