package ethereum

import (
	"context"
	"fmt"
	"math/big"
)

// SingleChainSender is a sender that can manage one chain
type SingleChainSender struct {
	ec *EthClient
}

// NewSingleChainSender creates a new SingleChainSender
func NewSingleChainSender(ec *EthClient) *SingleChainSender {
	return &SingleChainSender{ec}
}

// Send sends transaction
func (s *SingleChainSender) Send(chainID *big.Int, raw string) error {
	return s.ec.SendRawTransaction(context.Background(), raw)
}

// MultiChainSender is a sender that can manage one chain
type MultiChainSender struct {
	ecRegistry map[string]*EthClient
}

func chainIDToString(chainID *big.Int) string {
	return chainID.Text(16)
}

// NewMultiChainSender creates a new MultiChainSender
func NewMultiChainSender(ecs []*EthClient) *MultiChainSender {
	ecRegistry := make(map[string]*EthClient)
	for _, ec := range ecs {
		chainID, err := ec.NetworkID(context.Background())
		if err != nil {
			panic(err)
		}
		ecRegistry[chainIDToString(chainID)] = ec
	}
	return &MultiChainSender{ecRegistry}
}

// Send sends transaction
func (s *MultiChainSender) Send(chainID *big.Int, raw string) error {
	ec, ok := s.ecRegistry[chainIDToString(chainID)]
	if !ok {
		return fmt.Errorf("Not connected to chain %v", chainID)
	}
	return ec.SendRawTransaction(context.Background(), raw)
}
