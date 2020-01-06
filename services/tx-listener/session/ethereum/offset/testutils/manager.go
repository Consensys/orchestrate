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
	node := &dynamic.Node{
		ID: "test-node",
		Listener: &dynamic.Listener{
			BlockPosition: 10,
		},
	}

	lastBlockNumber, err := s.Manager.GetLastBlockNumber(context.Background(), node)
	assert.Nil(s.T(), err, "GetLastBlockNumber should not error")
	assert.Equal(s.T(), int64(0), lastBlockNumber, "Lastblock should be correct")

	err = s.Manager.SetLastBlockNumber(context.Background(), node, 12)
	assert.Nil(s.T(), err, "SetLastBlockNumber should not error")

	lastBlockNumber, err = s.Manager.GetLastBlockNumber(context.Background(), node)
	assert.Nil(s.T(), err, "GetLastBlockNumber should not error")
	assert.Equal(s.T(), int64(12), lastBlockNumber, "Lastblock should be correct")
}

func (s *OffsetManagerTestSuite) TestManagerLastIndex() {
	node := &dynamic.Node{
		ID: "test-node",
		Listener: &dynamic.Listener{
			BlockPosition: 10,
		},
	}

	lastTxIndex, err := s.Manager.GetLastTxIndex(context.Background(), node, 10)
	assert.Nil(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(0), lastTxIndex, "LastTxIndex should be correct")

	err = s.Manager.SetLastTxIndex(context.Background(), node, 10, 17)
	assert.Nil(s.T(), err, "SetLastTxIndex should not error")

	lastTxIndex, err = s.Manager.GetLastTxIndex(context.Background(), node, 10)
	assert.Nil(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(17), lastTxIndex, "LastTxIndex should be correct")

	lastTxIndex, err = s.Manager.GetLastTxIndex(context.Background(), node, 11)
	assert.Nil(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(0), lastTxIndex, "LastTxIndex should be correct")
}
