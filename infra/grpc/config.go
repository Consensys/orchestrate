package grpc

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	storeTargetFlag     = "grpc-store-target"
	storeTargetViperKey = "grpc.store.target"
	storeTargetDefault  = ""
	storeTargetEnv      = "GRPC_STORE_TARGET"
)

// StoreTarget register flag for Ethereum client URLs
func StoreTarget(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`GRPC Context Store target (See https://github.com/grpc/grpc/blob/master/doc/naming.md)
Environment variable: %q`, storeTargetEnv)
	f.String(storeTargetFlag, storeTargetDefault, desc)
	viper.SetDefault(storeTargetViperKey, storeTargetDefault)
	viper.BindPFlag(storeTargetViperKey, f.Lookup(storeTargetFlag))
	viper.BindEnv(storeTargetViperKey, storeTargetEnv)
}

// DialContext Create a new connection
func DialContext(ctx context.Context) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, viper.GetString(storeTargetViperKey))
}
