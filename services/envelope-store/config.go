package envelopestore

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(typeViperKey, typeDefault)
	_ = viper.BindEnv(typeViperKey, typeEnv)
}

var (
	typeFlag     = "envelope-store"
	typeViperKey = "envelope-store.type"
	typeDefault  = "pg"
	typeEnv      = "ENVELOPE_STORE"
)

// Type register flag for the Envelope Store to select
func Type(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Type of Envelope Store (one of %q)
Environment variable: %q`, []string{"mock", "pg"}, typeEnv)
	f.String(typeFlag, typeDefault, desc)
	_ = viper.BindPFlag(typeViperKey, f.Lookup(typeFlag))
}
