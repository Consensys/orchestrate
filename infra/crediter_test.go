package infra

import (
	"context"
	"math/big"
	"testing"

	"github.com/Shopify/sarama/mocks"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
)

func TestCrediter(t *testing.T) {
	mp := mocks.NewSyncProducer(t, nil)

	p, err := NewSaramaCrediter(mp)

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
