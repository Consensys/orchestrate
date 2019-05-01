package grpc

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGRPCStoreTarget(t *testing.T) {
	name := "grpc.store.target"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	StoreTarget(flgs)
	expected := ""
	assert.Equal(t, expected, viper.GetString(name), "Default")

	os.Setenv("GRPC_STORE_TARGET", "env-grpc-store")
	expected = "env-grpc-store"
	assert.Equal(t, expected, viper.GetString(name), "From Environment Variable")
	os.Unsetenv("GRPC_STORE_TARGET")

	args := []string{
		"--grpc-store-target=flag-grpc-store",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "Store should not error")
	expected = "flag-grpc-store"
	assert.Equal(t, expected, viper.GetString(name), "From Flag")
}
