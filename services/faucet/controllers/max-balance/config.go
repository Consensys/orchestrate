package maxbalance

import (
	"fmt"
	"math/big"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(faucetMaxViperKey, faucetMaxEnv)
	viper.SetDefault(faucetMaxViperKey, faucetMaxDefault)
}

var (
	faucetMaxFlag     = "faucet-max-balance"
	faucetMaxViperKey = "faucet.ctrl.max-balance"
	faucetMaxDefault  = "200000000000000000"
	faucetMaxEnv      = "FAUCET_MAX_BALANCE"
)

// FaucetMaxBalance register flag for Faucet Max Balance
func FaucetMaxBalance(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Faucet will stop crediting an address when it reaches this balance (Wei in decimal format)
Environment variable: %q`, faucetMaxEnv)
	f.String(faucetMaxFlag, faucetMaxDefault, desc)
	_ = viper.BindPFlag(faucetMaxViperKey, f.Lookup(faucetMaxFlag))
}

// Config holds MaxBalance Controller configuration
type Config struct {
	MaxBalance *big.Int
	BalanceAt  BalanceAtFunc
}
