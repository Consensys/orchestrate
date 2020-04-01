// +build unit

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAliasRegistry(t *testing.T) {
	reg := NewAliasRegistry()
	reg.Set("namespace-test-1", "foo", "test-foo")
	reg.Set("namespace-test-1", "bar", "test-bar-1")
	reg.Set("namespace-test-2", "bar", "test-bar-2")

	alias, _ := reg.Get("namespace-test-1", "foo")
	assert.Equal(t, "test-foo", alias, "#1 Alias should have been registered")
	alias, _ = reg.Get("namespace-test-1", "bar")
	assert.Equal(t, "test-bar-1", alias, "#2 Alias should have been registered")
	alias, _ = reg.Get("namespace-test-2", "bar")
	assert.Equal(t, "test-bar-2", alias, "#3 Alias should have been registered")
}
