package infra

import (
	"context"
	"math/big"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
	faucet "gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
)

// CreateFaucet create a faucet able to send a message to a kafka queue to credit an account with ethers
func CreateFaucet(ethURL string, faucetAddress common.Address, faucetMaxBalance int64, p sarama.SyncProducer, topic string) (*faucet.ControlledFaucet, error) {
	// Set BlackList controller
	chains := []*big.Int{
		big.NewInt(3), //
	}
	addresses := []common.Address{faucetAddress}
	bl := faucet.NewBlackList(chains, addresses)

	// Set Cooldown controller that requires a 60 sec interval between credits
	cd := faucet.NewCoolDown(time.Duration(60*time.Second), 50)

	// Set MaxBalance controller (it connects to ropsten to retrieve balance)
	ec, err := ethereum.Dial(ethURL)
	if err != nil {
		return &faucet.ControlledFaucet{}, err
	}
	balanceAt := func(chainID *big.Int, a common.Address) (*big.Int, error) {
		return ec.BalanceAt(context.Background(), a, nil)
	}
	mb := faucet.NewMaxBalance(
		big.NewInt(faucetMaxBalance),
		balanceAt,
	)

	// Create Faucet
	bc := baseCrediter(faucetAddress, p, topic)
	return faucet.NewControlledFaucet(bc, bl.Control, cd.Control, mb.Control), nil
}

func baseCrediter(faucetAddress common.Address, p sarama.SyncProducer, topic string) faucet.CreditFunc {
	return func(r *services.FaucetRequest) (*big.Int, bool, error) {
		msg, err := prepareFaucetMsg(r, faucetAddress)
		if err != nil {
			return r.Value, false, err
		}

		msg.Topic = topic
		_, _, err = p.SendMessage(&msg)
		if err != nil {
			return r.Value, false, err
		}

		return r.Value, true, nil
	}
}

func prepareFaucetMsg(r *services.FaucetRequest, faucetAddress common.Address) (sarama.ProducerMessage, error) {
	var msg sarama.ProducerMessage

	marshaller := infSarama.NewMarshaller()

	faucetTrace := types.NewTrace()
	faucetTrace.Chain().ID = r.ChainID
	faucetTrace.Sender().Address = &faucetAddress
	faucetTrace.Tx().SetValue(r.Value)
	faucetTrace.Tx().SetTo(&r.Address)

	err := marshaller.Marshal(faucetTrace, msg)
	if err != nil {
		return sarama.ProducerMessage{}, err
	}

	return msg, nil
}
