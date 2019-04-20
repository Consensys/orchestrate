package amount

import (
	"fmt"
	"math/big"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.BindEnv(faucetAmountViperKey, faucetAmountEnv)
	viper.SetDefault(faucetAmountViperKey, faucetAmountDefault)
}

var (
	faucetAmountFlag     = "faucet-amount"
	faucetAmountViperKey = "faucet.amount"
	faucetAmountDefault  = "100000000000000000"
	faucetAmountEnv      = "FAUCET_CREDIT_AMOUNT"
)

// FaucetAmount register flag for Faucet Amount
func FaucetAmount(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Amount to credit when calling Faucet (Wei in decimal format)
Environment variable: %q`, faucetAmountEnv)
	f.String(faucetAmountFlag, faucetAmountDefault, desc)
	viper.BindPFlag(faucetAmountViperKey, f.Lookup(faucetAmountFlag))
}

// Config for fixed amount controller
type Config struct {
	Amount *big.Int
}
