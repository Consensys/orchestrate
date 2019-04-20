package maxbalance

import (
	"fmt"
	"math/big"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.BindEnv(faucetMaxViperKey, faucetMaxEnv)
	viper.SetDefault(faucetMaxViperKey, faucetMaxDefault)
}

var (
	faucetMaxFlag     = "faucet-max"
	faucetMaxViperKey = "faucet.ctrl.max"
	faucetMaxDefault  = "200000000000000000"
	faucetMaxEnv      = "FAUCET_MAX_BALANCE"
)

// FaucetMaxBalance register flag for Faucet Max Balance
func FaucetMaxBalance(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Max balance (Wei in decimal format)
Environment variable: %q`, faucetMaxEnv)
	f.String(faucetMaxFlag, faucetMaxDefault, desc)
	viper.BindPFlag(faucetMaxViperKey, f.Lookup(faucetMaxFlag))
}

// Config holds MaxBalance Controller configuration
type Config struct {
	MaxBalance *big.Int
	BalanceAt  BalanceAtFunc
}
