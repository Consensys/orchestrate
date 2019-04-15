package infra

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	faucet "gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git"
)

func parseBlackList(blacklist []string) ([]*big.Int, []common.Address, error) {
	// Set BlackList controller
	var chains []*big.Int
	var addresses []common.Address
	for _, bl := range blacklist {
		split := strings.Split(bl, "@")
		chainID := big.NewInt(0)
		chainID, ok := chainID.SetString(split[1], 10)
		if !ok {
			return nil, nil, fmt.Errorf("Could not parse %v", bl)
		}
		chains = append(chains, chainID)
		addresses = append(addresses, common.HexToAddress(split[0]))
	}
	return chains, addresses, nil
}

func parseBalance(s string) (*big.Int, error) {
	balance := big.NewInt(0)
	balance, ok := balance.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("Could not parse max balance %q", s)
	}
	return balance, nil
}

// CreateFaucet create a faucet able to send a message to a kafka queue to credit an account with ethers
func CreateFaucet(balanceAt faucet.BalanceAtFunc, credit faucet.CreditFunc) (*faucet.ControlledFaucet, error) {

	chains, addresses, err := parseBlackList(viper.GetStringSlice("faucet.blacklist"))
	if err != nil {
		return nil, err
	}
	bl := faucet.NewBlackList(chains, addresses)

	// Set Cooldown controller that requires a 60 sec interval between credits
	cd := faucet.NewCoolDown(viper.GetDuration("faucet.cooldown"), 50)

	// Set MaxBalance controller
	maxBalance, err := parseBalance(viper.GetString("faucet.max"))
	if err != nil {
		return nil, err
	}

	mb := faucet.NewMaxBalance(
		maxBalance,
		balanceAt,
	)

	// Create Faucet
	return faucet.NewControlledFaucet(credit, bl.Control, cd.Control, mb.Control), nil
}
