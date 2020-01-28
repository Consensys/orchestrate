package testutils

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset"
)

// OffsetManagerTestSuite is a test suite for Offsets Manager
type OffsetManagerTestSuite struct {
	suite.Suite
	Manager offset.Manager
}

func (s *OffsetManagerTestSuite) TestManagerLastBlock() {
	chain := &dynamic.Chain{
		UUID: "test-chain",
		Listener: &dynamic.Listener{
			BlockPosition: 10,
		},
	}

	lastBlockNumber, err := s.Manager.GetLastBlockNumber(context.Background(), chain)
	assert.Nil(s.T(), err, "GetLastBlockNumber should not error")
	assert.Equal(s.T(), int64(0), lastBlockNumber, "Lastblock should be correct")

	err = s.Manager.SetLastBlockNumber(context.Background(), chain, 12)
	assert.Nil(s.T(), err, "SetLastBlockNumber should not error")

	lastBlockNumber, err = s.Manager.GetLastBlockNumber(context.Background(), chain)
	assert.Nil(s.T(), err, "GetLastBlockNumber should not error")
	assert.Equal(s.T(), int64(12), lastBlockNumber, "Lastblock should be correct")
}

func (s *OffsetManagerTestSuite) TestManagerLastIndex() {
	chain := &dynamic.Chain{
		UUID: "test-chain",
		Listener: &dynamic.Listener{
			BlockPosition: 10,
		},
	}

	lastTxIndex, err := s.Manager.GetLastTxIndex(context.Background(), chain, 10)
	assert.Nil(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(0), lastTxIndex, "LastTxIndex should be correct")

	err = s.Manager.SetLastTxIndex(context.Background(), chain, 10, 17)
	assert.Nil(s.T(), err, "SetLastTxIndex should not error")

	lastTxIndex, err = s.Manager.GetLastTxIndex(context.Background(), chain, 10)
	assert.Nil(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(17), lastTxIndex, "LastTxIndex should be correct")

	lastTxIndex, err = s.Manager.GetLastTxIndex(context.Background(), chain, 11)
	assert.Nil(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(0), lastTxIndex, "LastTxIndex should be correct")
}
