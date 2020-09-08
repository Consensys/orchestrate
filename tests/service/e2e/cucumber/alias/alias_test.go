// +build unit

package alias

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAliasRegistry(t *testing.T) {
	reg := NewAliasRegistry()
	reg.Set("test-foo", "namespace-test-1", "foo")
	reg.Set("test-bar-1", "namespace-test-1", "bar")
	reg.Set("test-bar-2", "namespace-test-2", "bar")

	alias, _ := reg.Get("namespace-test-1", "foo")
	assert.Equal(t, "test-foo", alias, "#1 Alias should have been registered")
	alias, _ = reg.Get("namespace-test-1", "bar")
	assert.Equal(t, "test-bar-1", alias, "#2 Alias should have been registered")
	alias, _ = reg.Get("namespace-test-2", "bar")
	assert.Equal(t, "test-bar-2", alias, "#3 Alias should have been registered")
}
