package loader

import (
	"context"
	"math/rand"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
	clientmock "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store/client/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

var letterRunes = []rune("abcdef0123456789")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func makeContext(i int, store *clientmock.EnvelopeStoreClient) *engine.TxContext {
	// Initialize context
	txctx := engine.NewTxContext().Prepare(log.NewEntry(log.StandardLogger()), nil)

	hash, uuid := ethereum.HexToHash("0x"+RandString(64)), RandString(32)

	// Prepare context
	txctx.Envelope.Chain = chain.FromInt(10)
	txctx.Envelope.Receipt = &ethereum.Receipt{
		TxHash: hash,
	}

	switch i % 2 {
	case 0:
		// Prestore an envelope
		evlp := &envelope.Envelope{
			Chain: chain.FromInt(10),
			Tx: &ethereum.Transaction{
				Hash: hash,
				Raw:  ethereum.HexToData("0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80"),
			},
			Metadata: &envelope.Metadata{Id: uuid},
		}
		_, _ = store.Store(
			context.Background(),
			&evlpstore.StoreRequest{
				Envelope: evlp,
			},
		)
		txctx.Set("uuid", uuid)
		txctx.Set("errors", 0)
		txctx.Set("status", "MINED")
	case 1:
		txctx.Set("uuid", "")
		txctx.Set("errors", 1)
		txctx.Set("status", "")
	}
	return txctx
}

type StoreTestSuite struct {
	testutils.HandlerTestSuite
	store *clientmock.EnvelopeStoreClient
}

func (s *StoreTestSuite) SetupSuite() {
	s.store = clientmock.New()
	s.Handler = EnvelopeLoader(s.store)
}

func (s *StoreTestSuite) TestStore() {
	rounds := 2
	var txctxs []*engine.TxContext
	for i := 0; i < rounds; i++ {
		txctxs = append(txctxs, makeContext(i, s.store))
	}

	// Handle contexts
	s.Handle(txctxs)

	for _, out := range txctxs {
		assert.Len(s.T(), out.Envelope.Errors, out.Get("errors").(int), "Expected correct errors count")
		assert.Equal(s.T(), out.Get("uuid").(string), out.Envelope.GetMetadata().GetId(), "Metadata should be set")
		if len(out.Envelope.Errors) == 0 {
			resp, _ := s.store.LoadByID(
				context.Background(),
				&evlpstore.LoadByIDRequest{
					Id: out.Envelope.GetMetadata().GetId(),
				},
			)
			assert.Equal(s.T(), out.Get("status").(string), resp.GetStatusInfo().GetStatus().String(), "Envelope should have expected status")
		}

	}

}

func TestLoader(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}
