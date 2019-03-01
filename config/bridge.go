package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	bridgeLinksFlag     = "bridge-links"
	bridgeLinksViperKey = "bridge.links"
	bridgeLinksEnv = "BRIDGE_LINKS"
)

// BridgeLinks lists bridges
func BridgeLinks(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`List of bridges with format addr@chainID<>addr@chainID: %q`, bridgeLinksEnv)
	f.StringSlice(bridgeLinksFlag, []string{}, desc)
	viper.BindPFlag(bridgeLinksViperKey, f.Lookup(bridgeLinksFlag))
	viper.BindEnv(bridgeLinksViperKey, bridgeLinksEnv)
}

var (
	bridgeMethodSignatureFlag     = "bridge-methodsignature"
	bridgeMethodSignatureViperKey = "bridge.methodsignature"
	bridgeMethodSignatureDefault  = "RelayMessage(bytes,address,address)"
	bridgeMethodSignatureEnv = "BRIDGE_METHODSIGNATURE"
)

// BridgeMethodSignature lists bridges
func BridgeMethodSignature(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Method signature to listen to when bridging. Format: methodName(typeParam1,typeParam2,...): %q`, bridgeMethodSignatureEnv)
	f.String(bridgeMethodSignatureFlag, bridgeMethodSignatureDefault, desc)
	viper.BindPFlag(bridgeMethodSignatureViperKey, f.Lookup(bridgeMethodSignatureFlag))
	viper.BindEnv(bridgeMethodSignatureViperKey, bridgeMethodSignatureEnv)
}

var (
	bridgeAuthorityFlag     = "bridge-authority"
	bridgeAuthorityViperKey = "bridge.authority"
	bridgeAuthorityEnv = "BRIDGE_AUTHORITY"
)

// BridgeAuthority lists bridges
func BridgeAuthority(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address used to sign the transaction on the destination chain.: %q`, bridgeAuthorityEnv)
	f.String(bridgeAuthorityFlag, "", desc)
	viper.BindPFlag(bridgeAuthorityViperKey, f.Lookup(bridgeAuthorityFlag))
	viper.BindEnv(bridgeAuthorityViperKey, bridgeAuthorityEnv)
}