package steps

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
)

func TestGetChainCounts(t *testing.T) {
	envelopes := map[string]*envelope.Envelope{
		"1": &envelope.Envelope{
			Chain: &common.Chain{Id: "888"},
		},
		"2": &envelope.Envelope{
			Chain: &common.Chain{Id: "777"},
		},
		"3": &envelope.Envelope{
			Chain: &common.Chain{Id: "888"},
		},
	}

	counts := GetChainCounts(envelopes)
	expected := map[string]uint{"777": 1, "888": 2}

	eq := reflect.DeepEqual(counts, expected)
	if !eq {
		t.Errorf("utils: error on counting")
	}
}

func TestChanTimeout(t *testing.T) {

	testChan := make(chan *envelope.Envelope)

	e, err := ChanTimeout(testChan, 1, 1)
	assert.Nil(t, e, "Should get an nil envelope slice")
	assert.Error(t, err, "Should get an error for not having received an envelope")

	testEnvelope := &envelope.Envelope{Chain: &common.Chain{Id: "888"}}
	go func() {
		testChan <- testEnvelope
	}()
	e, err = ChanTimeout(testChan, 1, 1)
	assert.NoError(t, err, "Should not get an error")
	assert.Equal(t, testEnvelope.GetChain().GetId(), e[0].GetChain().GetId(), "Should be the same envelope")
}
