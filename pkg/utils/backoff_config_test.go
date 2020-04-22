// +build unit

package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/spf13/viper"
)

func TestNewRetryConfig(t *testing.T) {
	r := NewRetryConfig(viper.New())
	assert.NotNil(t, r, "Should get a retry config")
}

func TestNewConfig(t *testing.T) {
	r := NewConfig(viper.New())
	assert.NotNil(t, r, "Should get a config")
}
