package multitenancy

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMultiTenancyEnable(t *testing.T) {
	name := "multi.tenancy.enabled"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Enabled(flgs)

	expected := false
	assert.Equal(t, expected, viper.GetBool(name), "TenancyEnable #1")
}

func TestTenantNamespace(t *testing.T) {
	name := "tenant.namespace"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Enabled(flgs)

	expected := false
	assert.Equal(t, expected, viper.GetBool(name), "TenantNamespace #1")
}
