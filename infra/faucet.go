package infra

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	faucet "gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git"
)

// CreateFaucet create a faucet able to send a message to a kafka queue to credit an account with ethers
func CreateFaucet(conf FaucetConfig, balanceAt faucet.BalanceAtFunc, credit faucet.CreditFunc) (*faucet.ControlledFaucet, error) {
	// Set BlackList controller
	var chains []*big.Int
	var addresses []common.Address
	for _, bl := range conf.BlackList {
		split := strings.Split(bl, "-")
		chainID := big.NewInt(0)
		chainID, ok := chainID.SetString(split[0], 10)
		if !ok {
			return nil, fmt.Errorf("Could not parse %v", bl)
		}
		chains = append(chains, chainID)
		addresses = append(addresses, common.HexToAddress(split[1]))
	}
	bl := faucet.NewBlackList(chains, addresses)

	// Set Cooldown controller that requires a 60 sec interval between credits
	cdDuration, err := time.ParseDuration(conf.CoolDownTime)
	if err != nil {
		return nil, err
	}
	cd := faucet.NewCoolDown(cdDuration, 50)

	// Set MaxBalance controller
	maxBalance := big.NewInt(0)
	maxBalance, ok := maxBalance.SetString(conf.MaxBalance, 10)
	if !ok {
		if err != nil {
			return nil, fmt.Errorf("Could not parse balance %q", conf.MaxBalance)
		}
	}
	mb := faucet.NewMaxBalance(
		maxBalance,
		balanceAt,
	)

	// Create Faucet
	return faucet.NewControlledFaucet(credit, bl.Control, cd.Control, mb.Control), nil
}
