package infra

import (
	"context"
	"math/big"
	"testing"

	"github.com/Shopify/sarama/mocks"
	"github.com/ethereum/go-ethereum/common"
	flags "github.com/jessevdk/go-flags"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/services"
)

func TestCrediter(t *testing.T) {
	mp := mocks.NewSyncProducer(t, nil)

	// Set default configuration
	opts := FaucetConfig{}
	flags.ParseArgs(&opts, []string{})
	p, err := NewSaramaCrediter(opts, mp)

	if err != nil {
		t.Errorf("Expected to create crediter with default config but got %v", err)
	}
	mp.ExpectSendMessageAndSucceed()

	amount, ok, err := p.Credit(
		context.Background(),
		&services.FaucetRequest{
			ChainID: big.NewInt(0),
			Value:   big.NewInt(13),
			Address: common.HexToAddress("0x664895b5fE3ddf049d2Fb508cfA03923859763C6"),
		},
	)

	if !ok || amount.Uint64() != 13 || err != nil {
		t.Errorf("Expected valid credit but got %v %v %v", amount, ok, err)
	}
}
